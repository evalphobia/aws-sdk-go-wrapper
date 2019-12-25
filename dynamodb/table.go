package dynamodb

import (
	"fmt"
	"strings"

	SDK "github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

const batchWriteItemMax = 25

// Table is a wapper struct for DynamoDB table
type Table struct {
	service        *DynamoDB
	name           string
	nameWithPrefix string
	design         *TableDesign

	putSpool   []*SDK.PutItemInput
	errorItems []*SDK.PutItemInput
}

// ---------------------------------
// table
// ---------------------------------

// NewTable returns initialized *Table.
func NewTable(svc *DynamoDB, name string) (*Table, error) {
	tableName := svc.prefix + name
	req, err := svc.client.DescribeTable(&SDK.DescribeTableInput{
		TableName: pointers.String(tableName),
	})
	if err != nil {
		svc.Errorf("error on `DescribeTable` operation; table=%s; error=%s;", name, err.Error())
		return nil, err
	}

	design := newTableDesignFromDescription(req.Table)
	return &Table{
		service:        svc,
		name:           name,
		nameWithPrefix: tableName,
		design:         design,
	}, nil
}

// NewTableWithDesign returns initialized *Table.
func NewTableWithDesign(svc *DynamoDB, design *TableDesign) (*Table, error) {
	tableName := design.name
	name := strings.Replace(tableName, svc.prefix, "", 1)
	return &Table{
		service:        svc,
		name:           name,
		nameWithPrefix: tableName,
		design:         design,
	}, nil
}

// NewTableWithoutDesign returns initialized *Table without table design.
func NewTableWithoutDesign(svc *DynamoDB, name string) *Table {
	tableName := svc.prefix + name
	return &Table{
		service:        svc,
		name:           name,
		nameWithPrefix: tableName,
	}
}

// GetDesign gets table design.
func (t *Table) GetDesign() *TableDesign {
	return t.design
}

// SetDesign sets table design.
func (t *Table) SetDesign(design *TableDesign) {
	t.design = design
}

// RefreshDesign returns refreshed table design.
func (t *Table) RefreshDesign() (*TableDesign, error) {
	req, err := t.service.client.DescribeTable(&SDK.DescribeTableInput{
		TableName: pointers.String(t.nameWithPrefix),
	})
	if err != nil {
		t.service.Errorf("error on `DescribeTable` operation; table=%s; error=%s", t.nameWithPrefix, err.Error())
		return nil, err
	}

	t.design = newTableDesignFromDescription(req.Table)
	return t.design, nil
}

// UpdateThroughput updates the r/w ProvisionedThroughput.
func (t *Table) UpdateThroughput(r int64, w int64) error {
	t.design.SetThroughput(r, w)
	return t.updateThroughput()
}

// UpdateWriteThroughput updates the write ProvisionedThroughput.
func (t *Table) UpdateWriteThroughput(w int64) error {
	t.design.SetThroughput(t.design.readCapacity, w)
	return t.updateThroughput()
}

// UpdateReadThroughput updates the read ProvisionedThroughput.
func (t *Table) UpdateReadThroughput(r int64) error {
	t.design.SetThroughput(r, t.design.writeCapacity)
	return t.updateThroughput()
}

// updateThroughput updates dynamodb table provisioned throughput
func (t *Table) updateThroughput() error {
	_, err := t.service.client.UpdateTable(&SDK.UpdateTableInput{
		TableName: pointers.String(t.nameWithPrefix),
		ProvisionedThroughput: &SDK.ProvisionedThroughput{
			ReadCapacityUnits:  pointers.Long64(t.design.readCapacity),
			WriteCapacityUnits: pointers.Long64(t.design.writeCapacity),
		},
	})
	if err != nil {
		t.service.Errorf("error on `UpdateTable` operation; table=%s; error=%s", t.nameWithPrefix, err.Error())
		return err
	}

	// refresh table information
	design, err := t.RefreshDesign()
	if err != nil {
		return err
	}
	t.design = design
	return nil
}

// ---------------------------------
// Put
// ---------------------------------

// AddItem adds an item to the write-waiting list (writeItem)
func (t *Table) AddItem(item *PutItem) {
	w := &SDK.PutItemInput{
		TableName:              pointers.String(t.nameWithPrefix),
		ReturnConsumedCapacity: pointers.String("TOTAL"),
		Item:                   item.data,
		Expected:               item.conditions,
	}
	t.putSpool = append(t.putSpool, w)
	t.service.addWriteTable(t)
}

// Put executes put operation from the write-waiting list (writeItem)
func (t *Table) Put() error {
	errList := newErrors()
	// アイテムの保存処理
	for _, item := range t.putSpool {
		err := t.validatePutItem(item)
		if err != nil {
			errList.Add(err)
			continue
		}

		_, err = t.service.client.PutItem(item)
		if err != nil {
			errList.Add(err)
			t.errorItems = append(t.errorItems, item)
		}
	}

	t.putSpool = nil
	if errList.HasError() {
		t.service.Errorf("errors on `Put` operations; table=%s; errors=[%s];", t.nameWithPrefix, errList.Error())
		return errList
	}
	return nil
}

// BatchPut executes BatchWriteItem operation from the write-waiting list (writeItem)
func (t *Table) BatchPut() error {
	errList := newErrors()
	errorSpoolIndices := make([]int, 0, len(t.putSpool))
	for index, item := range t.putSpool {
		err := t.validatePutItem(item)
		if err != nil {
			errList.Add(err)
			// バリデーションエラーのアイテムは送信スプールから除外するリストに加える
			errorSpoolIndices = append(errorSpoolIndices, index)
			continue
		}
	}
	t.removeErroredSpoolByIndices(errorSpoolIndices)
	input := new(SDK.BatchWriteItemInput)
	writeRequests := t.spoolToWriteRequests()
	for i := 0; i < len(writeRequests); i++ {
		input.SetRequestItems(writeRequests[i])
		if _, err := t.service.client.BatchWriteItem(input); err != nil {
			errList.Add(err)
		}
	}

	t.putSpool = nil
	if errList.HasError() {
		t.service.Errorf("errors on `Put` operations; table=%s; errors=[%s];", t.nameWithPrefix, errList.Error())
		return errList
	}
	return nil
}

// removeErroredSpoolByIndices removes elements which have index in validation error indices list from putSpool.
func (t *Table) removeErroredSpoolByIndices(errorSpoolIndices []int) {
	for i := len(errorSpoolIndices) - 1; i >= 0; i-- {
		removeIndex := errorSpoolIndices[i]
		firstHalf, latterHalf := t.putSpool[:removeIndex], t.putSpool[removeIndex+1:]
		t.putSpool = append(firstHalf, latterHalf...)
	}
}

// check if exists all primary keys in the item to write it.
func (t *Table) validatePutItem(item *SDK.PutItemInput) error {
	hashKey := t.design.GetHashKeyName()
	itemAttrs := item.Item
	if _, ok := itemAttrs[hashKey]; !ok {
		return fmt.Errorf("error on `validatePutItem`; No such hash key; table=%s; hashkey=%s", t.nameWithPrefix, hashKey)
	}

	rangeKey := t.design.GetRangeKeyName()
	if rangeKey == "" {
		return nil
	}

	if _, ok := itemAttrs[rangeKey]; !ok {
		return fmt.Errorf("error on `validatePutItem`; No such range key; table=%s; rangekey=%s", t.nameWithPrefix, rangeKey)
	}
	return nil
}

func (t *Table) spoolToWriteRequests() []map[string][]*SDK.WriteRequest {
	requestChunkCount := 1 + (len(t.putSpool) / batchWriteItemMax)
	writeRequestsChunks := make([]map[string][]*SDK.WriteRequest, 0, requestChunkCount)

	for chunkNumber := 0; chunkNumber < requestChunkCount; chunkNumber++ {
		offsetInSpool := batchWriteItemMax * chunkNumber
		if offsetInSpool >= len(t.putSpool) {
			break
		}
		writeRequests := make([]*SDK.WriteRequest, 0, batchWriteItemMax)
		for itemInChunk := 0; itemInChunk < batchWriteItemMax && offsetInSpool+itemInChunk < len(t.putSpool); itemInChunk++ {
			spoolIndex := batchWriteItemMax*chunkNumber + itemInChunk
			wr := new(SDK.WriteRequest)
			wr.SetPutRequest(&SDK.PutRequest{Item: t.putSpool[spoolIndex].Item})
			writeRequests = append(writeRequests, wr)
		}
		result := make(map[string][]*SDK.WriteRequest)
		result[t.nameWithPrefix] = writeRequests
		writeRequestsChunks = append(writeRequestsChunks, result)
	}

	return writeRequestsChunks
}

// ---------------------------------
// Scan
// ---------------------------------

// Scan executes Scan operation.
func (t *Table) Scan() (*QueryResult, error) {
	cond := t.NewConditionList()
	cond.SetLimit(1000)
	return t.scan(cond, &SDK.ScanInput{})
}

// ScanWithCondition executes Scan operation with given condition.
func (t *Table) ScanWithCondition(cond *ConditionList) (*QueryResult, error) {
	return t.scan(cond, &SDK.ScanInput{})
}

// scan executes Scan operation.
func (t *Table) scan(cond *ConditionList, in *SDK.ScanInput) (*QueryResult, error) {
	if cond.HasFilter() {
		in.FilterExpression = cond.FormatFilter()
		in.ExpressionAttributeValues = cond.FormatValues()
		in.ExpressionAttributeNames = cond.FormatNames()
	}

	if cond.HasIndex() {
		in.IndexName = pointers.String(cond.index)
	}
	if cond.HasLimit() {
		in.Limit = pointers.Long64(cond.limit)
	}
	if cond.isConsistent {
		in.ConsistentRead = pointers.Bool(cond.isConsistent)
	}

	in.ExclusiveStartKey = cond.startKey
	in.TableName = pointers.String(t.nameWithPrefix)
	req, err := t.service.client.Scan(in)
	if err != nil {
		t.service.Errorf("error on `Scan` operation; table=%s; error=%s;", t.nameWithPrefix, err.Error())
		return nil, err
	}

	res := &QueryResult{
		Items:            req.Items,
		LastEvaluatedKey: req.LastEvaluatedKey,
		Count:            *req.Count,
		ScannedCount:     *req.ScannedCount,
	}
	return res, nil
}

// ---------------------------------
// Query
// ---------------------------------

// Query executes Query operation.
func (t *Table) Query(cond *ConditionList) (*QueryResult, error) {
	return t.query(cond, &SDK.QueryInput{})
}

// Count executes Query operation and get Count.
func (t *Table) Count(cond *ConditionList) (*QueryResult, error) {
	return t.query(cond, &SDK.QueryInput{
		Select: pointers.String("COUNT"),
	})
}

func (t *Table) query(cond *ConditionList, in *SDK.QueryInput) (*QueryResult, error) {
	if !cond.HasCondition() {
		err := fmt.Errorf("condition is missing, you must specify at least one condition")
		t.service.Errorf("error on `query`; table=%s; error=%s", t.nameWithPrefix, err.Error())
		return nil, err
	}

	in.KeyConditionExpression = cond.FormatCondition()
	in.FilterExpression = cond.FormatFilter()
	in.ExpressionAttributeValues = cond.FormatValues()
	in.ExpressionAttributeNames = cond.FormatNames()

	if cond.HasIndex() {
		in.IndexName = pointers.String(cond.index)
	}
	if cond.HasLimit() {
		in.Limit = pointers.Long64(cond.limit)
	}
	if cond.isConsistent {
		in.ConsistentRead = pointers.Bool(cond.isConsistent)
	}
	if cond.isDesc {
		in.ScanIndexForward = pointers.Bool(false)
	}

	in.ExclusiveStartKey = cond.startKey
	in.TableName = pointers.String(t.nameWithPrefix)
	req, err := t.service.client.Query(in)
	if err != nil {
		t.service.Errorf("error on `Query` operation; table=%s; error=%s", t.nameWithPrefix, err.Error())
		return nil, err
	}

	res := &QueryResult{
		Items:            req.Items,
		LastEvaluatedKey: req.LastEvaluatedKey,
		Count:            *req.Count,
		ScannedCount:     *req.ScannedCount,
	}
	return res, nil
}

// NewConditionList returns initialized *ConditionList.
func (t *Table) NewConditionList() *ConditionList {
	return NewConditionList(t.design.GetKeyAttributes())
}

// ---------------------------------
// Get
// ---------------------------------

// GetOne retrieves a single item by GetOne(HashKey [, RangeKey])
func (t *Table) GetOne(hashValue interface{}, rangeValue ...interface{}) (map[string]interface{}, error) {
	in := &SDK.GetItemInput{
		TableName: pointers.String(t.nameWithPrefix),
		Key:       t.design.keyAttributeValue(hashValue, rangeValue...),
	}
	req, err := t.service.client.GetItem(in)
	switch {
	case err != nil:
		t.service.Errorf("error on `GetItem` operation; table=%s; error=%s", t.nameWithPrefix, err.Error())
		return nil, err
	case req.Item == nil:
		return nil, nil
	}

	return UnmarshalAttributeValue(req.Item), nil
}

// ---------------------------------
// Delete
// ---------------------------------

// Delete deletes the item.
func (t *Table) Delete(hashValue interface{}, rangeValue ...interface{}) error {
	in := &SDK.DeleteItemInput{
		TableName: pointers.String(t.nameWithPrefix),
		Key:       t.design.keyAttributeValue(hashValue, rangeValue...),
	}

	_, err := t.service.client.DeleteItem(in)
	if err != nil {
		t.service.Errorf("error on `DeleteItem` operation; table=%s; error=%s", t.nameWithPrefix, err.Error())
		return err
	}
	return nil
}

// ForceDeleteAll deltes all data in the table.
// This performs scan all item and delete it each one by one.
func (t *Table) ForceDeleteAll() error {
	hashkey := t.design.GetHashKeyName()
	rangekey := t.design.GetRangeKeyName()

	result, err := t.Scan()
	if err != nil {
		return err
	}

	errData := newErrors()
	for _, item := range result.ToSliceMap() {
		var e error
		switch rangekey {
		case "":
			e = t.Delete(item[hashkey])
		default:
			e = t.Delete(item[hashkey], item[rangekey])
		}

		if e != nil {
			errData.Add(e)
		}
	}

	if errData.HasError() {
		return errData
	}
	return nil
}
