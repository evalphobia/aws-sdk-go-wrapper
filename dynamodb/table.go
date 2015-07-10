// DynamoDB Table operation/manipuration

package dynamodb

import (
	SDK "github.com/aws/aws-sdk-go/service/dynamodb"

	"errors"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
	"strings"
)

// DynamoTable is a wapper struct for DynamoDB table
type DynamoTable struct {
	db         *AmazonDynamoDB
	table      *SDK.TableDescription
	name       string
	indexes    map[string]*DynamoIndex
	writeItems []*SDK.PutItemInput
	errorItems []*SDK.PutItemInput
}

func (t *DynamoTable) Desc() (*SDK.TableDescription, error) {
	req, err := t.db.client.DescribeTable(&SDK.DescribeTableInput{
		TableName: String(t.name),
	})
	if err != nil {
		log.Error("[DynamoDB] Error in `DescribeTable` operation, table="+t.name, err)
		return nil, err
	}
	return req.Table, nil
}

func (t *DynamoTable) UpdateThroughput(r int64, w int64) error {
	th := t.table.ProvisionedThroughput
	th.ReadCapacityUnits = Long(r)
	th.WriteCapacityUnits = Long(w)
	return t.updateThroughput(t.toProvisionedThroughput(th))
}

func (t *DynamoTable) UpdateWriteThroughput(v int64) error {
	th := t.table.ProvisionedThroughput
	th.WriteCapacityUnits = Long(v)
	return t.updateThroughput(t.toProvisionedThroughput(th))
}

func (t *DynamoTable) UpdateReadThroughput(v int64) error {
	th := t.table.ProvisionedThroughput
	th.ReadCapacityUnits = Long(v)
	return t.updateThroughput(t.toProvisionedThroughput(th))
}

func (t *DynamoTable) toProvisionedThroughput(in *SDK.ProvisionedThroughputDescription) *SDK.UpdateTableInput {
	return &SDK.UpdateTableInput{
		TableName: String(t.name),
		ProvisionedThroughput: &SDK.ProvisionedThroughput{
			ReadCapacityUnits:  in.ReadCapacityUnits,
			WriteCapacityUnits: in.WriteCapacityUnits,
		},
	}
}

// updateThroughput updates dynamodb table provisioned throughput
func (t *DynamoTable) updateThroughput(in *SDK.UpdateTableInput) error {
	if t.isSameThroughput(in) {
		return nil
	}
	_, err := t.db.client.UpdateTable(in)
	if err != nil {
		log.Error("[DynamoDB] Error in `UpdateTable` operation, table="+t.name, err)
		return err
	}
	desc, err := t.Desc()
	if err != nil {
		return err
	}
	t.table = desc
	return nil
}

// isSameThroughput checks if the given throughput is same as current table throughput or not
func (t *DynamoTable) isSameThroughput(in *SDK.UpdateTableInput) bool {
	desc, err := t.Desc()
	if err != nil {
		return false
	}
	from := desc.ProvisionedThroughput
	to := in.ProvisionedThroughput
	switch {
	case *from.ReadCapacityUnits != *to.ReadCapacityUnits:
		return false
	case *from.WriteCapacityUnits != *to.WriteCapacityUnits:
		return false
	}
	return true
}

// AddItem adds an item to the write-waiting list (writeItem)
func (t *DynamoTable) AddItem(item *DynamoItem) {
	w := &SDK.PutItemInput{}
	w.TableName = String(t.name)
	w.ReturnConsumedCapacity = String("TOTAL")
	w.Item = item.data
	w.Expected = item.conditions
	t.writeItems = append(t.writeItems, w)
	t.db.addWriteTable(t.name)
}

// excecute write operation in the write-waiting list (writeItem)
func (t *DynamoTable) Put() error {
	var err error = nil
	var errs []string
	// アイテムの保存処理
	for _, item := range t.writeItems {
		if !t.isExistPrimaryKeys(item) {
			msg := "[DynamoDB] Cannot find primary key, table=" + t.name
			errs = append(errs, msg)
			log.Error(msg, item)
			continue
		}
		_, e := t.db.client.PutItem(item)
		if e != nil {
			errs = append(errs, e.Error())
			t.errorItems = append(t.errorItems, item)
		}
	}
	t.writeItems = []*SDK.PutItemInput{}
	if len(errs) != 0 {
		err = errors.New(strings.Join(errs, "\n"))
	}
	return err
}

// GetOne retrieves a single item by GetOne(HashKey [, RangeKey])
func (t *DynamoTable) GetOne(values ...Any) (map[string]interface{}, error) {
	key := NewItem()
	key.AddAttribute(t.GetHashKeyName(), values[0])
	if len(values) > 1 && t.GetRangeKeyName() != "" {
		key.AddAttribute(t.GetRangeKeyName(), values[1])
	}

	in := &SDK.GetItemInput{
		TableName: String(t.name),
		Key:       key.data,
	}
	req, err := t.db.client.GetItem(in)
	if err != nil {
		log.Error("[DynamoDB] Error in `GetItem` operation, table="+t.name, err)
		return nil, err
	}
	return Unmarshal(req.Item), nil
}

// query using LSI or GSI
func (t *DynamoTable) GetByIndex(idx string, values ...Any) ([]map[string]interface{}, error) {
	index, ok := t.indexes[idx]
	if !ok {
		log.Error("[DynamoDB] Cannot find the index name, table="+t.name, idx)
		log.Error("[DynamoDB] indexes on table="+t.name, t.indexes)
	}

	hashKey := index.GetHashKeyName()
	rangeKey := index.GetRangeKeyName()

	keys := make(map[string]*SDK.Condition)
	keys[hashKey] = &SDK.Condition{
		AttributeValueList: []*SDK.AttributeValue{createAttributeValue(values[0])},
		ComparisonOperator: String(ComparisonOperatorEQ),
	}
	if len(values) > 1 && rangeKey != "" {
		keys[rangeKey] = &SDK.Condition{
			AttributeValueList: []*SDK.AttributeValue{createAttributeValue(values[1])},
			ComparisonOperator: String(ComparisonOperatorEQ),
		}
	}

	in := &SDK.QueryInput{
		TableName:     String(t.name),
		KeyConditions: keys,
		IndexName:     &idx,
	}
	return t.Query(in)
}

// perform GetOne() and return slice value with single item
func (t *DynamoTable) getOneAsSlice(values []Any) ([]map[string]interface{}, error) {
	var (
		items []map[string]interface{}
		item  map[string]interface{}
		err   error
	)
	if len(values) > 1 {
		item, err = t.GetOne(values[0], values[1])
	} else {
		item, err = t.GetOne(values[0])
	}

	if err != nil {
		return items, err
	}
	if len(item) == 0 {
		return items, nil
	}
	return append(items, item), nil
}

// perform GetOne() or Query()
func (t *DynamoTable) Get(values ...Any) ([]map[string]interface{}, error) {
	if len(values) > 1 || t.GetRangeKeyName() == "" {
		return t.getOneAsSlice(values)
	}

	keys := make(map[string]*SDK.Condition)
	keys[t.GetHashKeyName()] = &SDK.Condition{
		AttributeValueList: []*SDK.AttributeValue{createAttributeValue(values[0])},
		ComparisonOperator: String(ComparisonOperatorEQ),
	}

	in := &SDK.QueryInput{
		TableName:     String(t.name),
		KeyConditions: keys,
	}
	return t.Query(in)
}

// get mapped-items with Query operation
func (t *DynamoTable) Query(in *SDK.QueryInput) ([]map[string]interface{}, error) {
	req, err := t.db.client.Query(in)
	if err != nil {
		log.Error("[DynamoDB] Error in `Query` operation, table="+t.name, err)
		return nil, err
	}
	return t.convertItemsToMapArray(req.Items), nil
}

// get mapped-items with Scan operation
func (t *DynamoTable) Scan() ([]map[string]interface{}, error) {
	in := &SDK.ScanInput{
		TableName: String(t.name),
		Limit:     Long(1000),
	}
	req, err := t.db.client.Scan(in)
	if err != nil {
		log.Error("[DynamoDB] Error in `Scan` operation, table="+t.name, err)
		return nil, err
	}
	return t.convertItemsToMapArray(req.Items), nil
}

// delete item
func (t *DynamoTable) Delete(values ...Any) error {
	key := NewItem()
	key.AddAttribute(t.GetHashKeyName(), values[0])
	if len(values) > 1 && t.GetRangeKeyName() != "" {
		key.AddAttribute(t.GetRangeKeyName(), values[1])
	}

	in := &SDK.DeleteItemInput{
		TableName: String(t.name),
		Key:       key.data,
	}
	_, err := t.db.client.DeleteItem(in)
	if err != nil {
		log.Error("[DynamoDB] Error in `DeleteItem` operation, table="+t.name, err)
		return err
	}
	return nil
}

// convert from dynamodb values to map
func (t *DynamoTable) convertItemsToMapArray(items []map[string]*SDK.AttributeValue) []map[string]interface{} {
	var m []map[string]interface{}
	for _, item := range items {
		m = append(m, Unmarshal(item))
	}
	return m
}

// get the name of hash key
func (t *DynamoTable) GetHashKeyName() string {
	return *t.table.KeySchema[0].AttributeName
}

// get the name of range key if exist
func (t *DynamoTable) GetRangeKeyName() string {
	if len(t.table.KeySchema) > 1 {
		return *t.table.KeySchema[1].AttributeName
	} else {
		return ""
	}
}

// check if exists all primary keys in the item to write it.
func (t *DynamoTable) isExistPrimaryKeys(item *SDK.PutItemInput) bool {
	hashKey := t.GetHashKeyName()
	itemAttrs := item.Item
	_, ok := itemAttrs[hashKey]
	if !ok {
		log.Warn("[DynamoDB] No HashKey, table="+t.name, hashKey)
		return false
	}
	rangeKey := t.GetRangeKeyName()
	_, ok = itemAttrs[rangeKey]
	if rangeKey != "" && !ok {
		log.Warn("[DynamoDB] No RangeKey, table="+t.name, rangeKey)
		return false
	}
	return true
}

// [CAUTION]
// only used this for developing, this performs scan all item and delete it each one by one
func (t *DynamoTable) DeleteAll() error {
	hashkey := t.GetHashKeyName()
	rangekey := t.GetRangeKeyName()

	result, err := t.Scan()
	if err != nil {
		return err
	}

	errStr := ""
	for _, item := range result {
		var e error
		if rangekey != "" {
			e = t.Delete(item[hashkey], item[rangekey])
		} else {
			e = t.Delete(item[hashkey])
		}
		if e != nil {
			errStr = errStr + "," + e.Error()
		}
	}

	if errStr != "" {
		err = errors.New(errStr)
	}

	return err
}
