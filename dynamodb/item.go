package dynamodb

import (
	SDK "github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

// PutItem is wrapped struct for DynamoDB Item to put.
type PutItem struct {
	data       map[string]*SDK.AttributeValue
	conditions map[string]*SDK.ExpectedAttributeValue
}

// NewPutItem returns initialized *PutItem.
func NewPutItem() *PutItem {
	return &PutItem{
		data:       make(map[string]*SDK.AttributeValue),
		conditions: make(map[string]*SDK.ExpectedAttributeValue),
	}
}

// AddAttribute adds an attribute to the PutItem.
func (item *PutItem) AddAttribute(name string, value interface{}) {
	item.data[name] = createAttributeValue(value)
}

// GetAttribute gets an attribute from PutItem.
func (item *PutItem) GetAttribute(name string) interface{} {
	return item.data[name]
}

// AddConditionExist adds a EXIST condition.
func (item *PutItem) AddConditionExist(name string) {
	item.addCondition(name, &SDK.ExpectedAttributeValue{
		Exists: pointers.Bool(true),
	})
}

// AddConditionNotExist adds a NOT EXIST condition.
func (item *PutItem) AddConditionNotExist(name string) {
	item.addCondition(name, &SDK.ExpectedAttributeValue{
		Exists: pointers.Bool(false),
	})
}

// AddConditionEQ adds a EQUAL condition.
func (item *PutItem) AddConditionEQ(name string, value interface{}) {
	cond := NewExpectedCondition(value, ComparisonOperatorEQ)
	item.addCondition(name, cond)
}

// AddConditionNE adds a NOT EQUAL condition.
func (item *PutItem) AddConditionNE(name string, value interface{}) {
	cond := NewExpectedCondition(value, ComparisonOperatorNE)
	item.addCondition(name, cond)
}

// AddConditionGT adds a GREATER THAN condition.
func (item *PutItem) AddConditionGT(name string, value interface{}) {
	cond := NewExpectedCondition(value, ComparisonOperatorGT)
	item.addCondition(name, cond)
}

// AddConditionLT adds a LESS THAN condition.
func (item *PutItem) AddConditionLT(name string, value interface{}) {
	cond := NewExpectedCondition(value, ComparisonOperatorLT)
	item.addCondition(name, cond)
}

// AddConditionGE adds a GREATER THAN or EQUAL condition.
func (item *PutItem) AddConditionGE(name string, value interface{}) {
	cond := NewExpectedCondition(value, ComparisonOperatorGE)
	item.addCondition(name, cond)
}

// AddConditionLE adds a LESS THAN or EQUAL condition.
func (item *PutItem) AddConditionLE(name string, value interface{}) {
	cond := NewExpectedCondition(value, ComparisonOperatorLE)
	item.addCondition(name, cond)
}

// addCondition adds a condition.
func (item *PutItem) addCondition(name string, condition *SDK.ExpectedAttributeValue) {
	item.conditions[name] = condition
}

// CountUp counts up the value.
func (item *PutItem) CountUp(name string, num int) {
	// TODO: implement atomic counter
}

// CountDown counts down the value.
func (item *PutItem) CountDown(name string, num int) {
	// TODO: implement atomic counter
}

// NewExpectedCondition returns *SDK.ExpectedAttributeValue with ComparisonOperator and value.
func NewExpectedCondition(value interface{}, operator string) *SDK.ExpectedAttributeValue {
	v := createAttributeValue(value)
	return &SDK.ExpectedAttributeValue{
		Value:              v,
		ComparisonOperator: pointers.String(operator),
	}
}
