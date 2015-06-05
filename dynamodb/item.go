// DynamoDB Item operation/manipuration

package dynamodb

import (
	SDK "github.com/aws/aws-sdk-go/service/dynamodb"
)

// wrapped struct for DynamoDB Item (data of each row in RDB concept)
type DynamoItem struct {
	data       map[string]*SDK.AttributeValue
	conditions map[string]*SDK.ExpectedAttributeValue
}

// Create new empty Item
func NewItem() *DynamoItem {
	return &DynamoItem{
		data:       make(map[string]*SDK.AttributeValue),
		conditions: make(map[string]*SDK.ExpectedAttributeValue),
	}
}

// Add a attribute to the Item
func (item *DynamoItem) AddAttribute(name string, value Any) {
	item.data[name] = createAttributeValue(value)
}

// Add a EXIST condition for put
func (item *DynamoItem) AddConditionExist(name string) {
	cond := NewExpected()
	cond.Exists = Boolean(true)
	item.AddCondition(name, cond)
}

// Add a NOT EXIST condition for put
func (item *DynamoItem) AddConditionNotExist(name string) {
	cond := NewExpected()
	cond.Exists = Boolean(false)
	item.AddCondition(name, cond)
}

// Add a EQUAL condition for put
func (item *DynamoItem) AddConditionEQ(name string, value Any) {
	cond := NewExpectedCondition(value, ComparisonOperatorEQ)
	item.AddCondition(name, cond)
}

// Add a NOT EQUAL condition for put
func (item *DynamoItem) AddConditionNE(name string, value Any) {
	cond := NewExpectedCondition(value, ComparisonOperatorNE)
	item.AddCondition(name, cond)
}

// Add a GREATER THAN condition for put
func (item *DynamoItem) AddConditionGT(name string, value Any) {
	cond := NewExpectedCondition(value, ComparisonOperatorGT)
	item.AddCondition(name, cond)
}

// Add a LESS THAN condition for put
func (item *DynamoItem) AddConditionLT(name string, value Any) {
	cond := NewExpectedCondition(value, ComparisonOperatorLT)
	item.AddCondition(name, cond)
}

// Add a GREATER THAN or EQUAL condition for put
func (item *DynamoItem) AddConditionGE(name string, value Any) {
	cond := NewExpectedCondition(value, ComparisonOperatorGE)
	item.AddCondition(name, cond)
}

// Add a LESS THAN or EQUAL condition for put
func (item *DynamoItem) AddConditionLE(name string, value Any) {
	cond := NewExpectedCondition(value, ComparisonOperatorLE)
	item.AddCondition(name, cond)
}

// Add a condition for put
func (item *DynamoItem) AddCondition(name string, condition *SDK.ExpectedAttributeValue) {
	item.conditions[name] = condition
}

// Atomic Counter
func (item *DynamoItem) CountUp(name string, num int) {
	// TODO: implement atomic counter
}

// Atomic Counter
func (item *DynamoItem) CountDown(name string, num int) {
	// TODO: implement atomic counter
}

// Create new empty condition for put
func NewExpected() *SDK.ExpectedAttributeValue {
	return &SDK.ExpectedAttributeValue{}
}

// Create new condition for put
func NewExpectedCondition(value Any, operator string) *SDK.ExpectedAttributeValue {
	v := createAttributeValue(value)
	return &SDK.ExpectedAttributeValue{
		Value:              v,
		ComparisonOperator: String(operator),
	}
}
