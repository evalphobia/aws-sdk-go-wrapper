// DynamoDB Item operation/manipuration

package dynamodb

import (
	AWS "github.com/awslabs/aws-sdk-go/aws"
	DynamoDB "github.com/awslabs/aws-sdk-go/gen/dynamodb"
)

// wrapped struct for DynamoDB Item (data of each row in RDB concept)
type DynamoItem struct {
	data       map[string]DynamoDB.AttributeValue
	conditions map[string]DynamoDB.ExpectedAttributeValue
}

// Create new empty Item
func NewItem() *DynamoItem {
	return &DynamoItem{
		data:       make(map[string]DynamoDB.AttributeValue),
		conditions: make(map[string]DynamoDB.ExpectedAttributeValue),
	}
}

// Add a attribute to the Item
func (item *DynamoItem) AddAttribute(name string, value Any) {
	item.data[name] = createAttributeValue(value)
}

// Add a EXIST condition for put
func (item *DynamoItem) AddConditionExist(name string) {
	cond := NewExpected()
	cond.Exists = AWS.Boolean(true)
	item.AddCondition(name, *cond)
}

// Add a NOT EXIST condition for put
func (item *DynamoItem) AddConditionNotExist(name string) {
	cond := NewExpected()
	cond.Exists = AWS.Boolean(false)
	item.AddCondition(name, *cond)
}

// Add a EQUAL condition for put
func (item *DynamoItem) AddConditionEQ(name string, value Any) {
	cond := NewExpectedCondition(value, DynamoDB.ComparisonOperatorEq)
	item.AddCondition(name, *cond)
}

// Add a NOT EQUAL condition for put
func (item *DynamoItem) AddConditionNE(name string, value Any) {
	cond := NewExpectedCondition(value, DynamoDB.ComparisonOperatorNe)
	item.AddCondition(name, *cond)
}

// Add a GREATER THAN condition for put
func (item *DynamoItem) AddConditionGT(name string, value Any) {
	cond := NewExpectedCondition(value, DynamoDB.ComparisonOperatorGt)
	item.AddCondition(name, *cond)
}

// Add a LESS THAN condition for put
func (item *DynamoItem) AddConditionLT(name string, value Any) {
	cond := NewExpectedCondition(value, DynamoDB.ComparisonOperatorLt)
	item.AddCondition(name, *cond)
}

// Add a GREATER THAN or EQUAL condition for put
func (item *DynamoItem) AddConditionGE(name string, value Any) {
	cond := NewExpectedCondition(value, DynamoDB.ComparisonOperatorGe)
	item.AddCondition(name, *cond)
}

// Add a LESS THAN or EQUAL condition for put
func (item *DynamoItem) AddConditionLE(name string, value Any) {
	cond := NewExpectedCondition(value, DynamoDB.ComparisonOperatorLe)
	item.AddCondition(name, *cond)
}

// Add a condition for put
func (item *DynamoItem) AddCondition(name string, condition DynamoDB.ExpectedAttributeValue) {
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
func NewExpected() *DynamoDB.ExpectedAttributeValue {
	return &DynamoDB.ExpectedAttributeValue{}
}

// Create new condition for put
func NewExpectedCondition(value Any, operator string) *DynamoDB.ExpectedAttributeValue {
	v := createAttributeValue(value)
	return &DynamoDB.ExpectedAttributeValue{
		Value:              &v,
		ComparisonOperator: AWS.String(operator),
	}
}
