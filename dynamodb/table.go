// DynamoDB Table operation/manipuration

package dynamodb

import (
	AWS "github.com/awslabs/aws-sdk-go/aws"
	DynamoDB "github.com/awslabs/aws-sdk-go/service/dynamodb"

	"errors"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

// DynamoDB table wrapper struct
type DynamoTable struct {
	db         *AmazonDynamoDB
	table      *DynamoDB.TableDescription
	name       string
	writeItems []*DynamoDB.PutItemInput
	errorItems []*DynamoDB.PutItemInput
}

// add item to the write-waiting list (writeItem)
func (t *DynamoTable) AddItem(item *DynamoItem) {
	w := &DynamoDB.PutItemInput{}
	w.TableName = AWS.String(t.name)
	w.ReturnConsumedCapacity = AWS.String("TOTAL")
	w.Item = item.data
	w.Expected = item.conditions
	t.writeItems = append(t.writeItems, w)
	t.db.addWriteTable(t.name)
}

// excecute write operation in the write-waiting list (writeItem)
func (t *DynamoTable) Put() error {
	var err error = nil
	errStr := ""
	// アイテムの保存処理
	for _, item := range t.writeItems {
		if !t.isExistPrimaryKeys(item) {
			log.Error("[DynamoDB] Cannot find primary key, table="+t.name, item)
			continue
		}
		_, e := t.db.client.PutItem(item)
		if e != nil {
			errStr = errStr + "," + e.Error()
			t.errorItems = append(t.errorItems, item)
		}
	}
	if errStr != "" {
		err = errors.New(errStr)
	}
	return err
}

// retrieve a single item
func (t *DynamoTable) GetOne(values ...Any) (map[string]interface{}, error) {
	key := NewItem()
	key.AddAttribute(t.GetHashKeyName(), values[0])
	if len(values) > 1 && t.GetRangeKeyName() != "" {
		key.AddAttribute(t.GetRangeKeyName(), values[1])
	}

	in := &DynamoDB.GetItemInput{
		TableName: AWS.String(t.name),
		Key:       key.data,
	}
	req, err := t.db.client.GetItem(in)
	if err != nil {
		log.Error("[DynamoDB] Error in `GetItem` operation, table="+t.name, err)
		return nil, err
	}
	return Unmarshal(req.Item), nil
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
	return append(items, item), nil
}

// perform GetOne() or Query()
func (t *DynamoTable) Get(values ...Any) ([]map[string]interface{}, error) {
	if len(values) > 1 || t.GetRangeKeyName() == "" {
		return t.getOneAsSlice(values)
	}

	keys := make(map[string]DynamoDB.Condition)
	keys[t.GetHashKeyName()] = DynamoDB.Condition{
		AttributeValueList: []DynamoDB.AttributeValue{createAttributeValue(values[0])},
		ComparisonOperator: AWS.String(DynamoDB.ComparisonOperatorEq),
	}

	in := &DynamoDB.QueryInput{
		TableName:     AWS.String(t.name),
		KeyConditions: keys,
	}
	return t.Query(in)
}

// get mapped-items with Query operation
func (t *DynamoTable) Query(in *DynamoDB.QueryInput) ([]map[string]interface{}, error) {
	req, err := t.db.client.Query(in)
	if err != nil {
		log.Error("[DynamoDB] Error in `Query` operation, table="+t.name, err)
		return nil, err
	}
	return t.convertItemsToMapArray(req.Items), nil
}

// get mapped-items with Scan operation
func (t *DynamoTable) Scan() ([]map[string]interface{}, error) {
	in := &DynamoDB.ScanInput{
		TableName: AWS.String(t.name),
		Limit:     AWS.Integer(1000),
	}
	req, err := t.db.client.Scan(in)
	if err != nil {
		log.Error("[DynamoDB] Error in `Scan` operation, table="+t.name, err)
		return nil, err
	}
	return t.convertItemsToMapArray(req.Items), nil
}

// convert from dynamodb values to map
func (t *DynamoTable) convertItemsToMapArray(items []map[string]DynamoDB.AttributeValue) []map[string]interface{} {
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
func (t *DynamoTable) isExistPrimaryKeys(item *DynamoDB.PutItemInput) bool {
	hashKey := t.GetHashKeyName()
	_, ok := item.Item[hashKey]
	if !ok {
		log.Warn("[DynamoDB] No HashKey, table="+t.name, hashKey)
		return false
	}
	rangeKey := t.GetRangeKeyName()
	_, ok = item.Item[rangeKey]
	if rangeKey != "" && !ok {
		log.Warn("[DynamoDB] No RangeKey, table="+t.name, rangeKey)
		return false
	}
	return true
}
