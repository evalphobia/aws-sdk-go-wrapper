// DynamoDB QueryCondition operation/manipuration

package dynamodb

import (
	SDK "github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	ComparisonOperatorEQ = "EQ"
	ComparisonOperatorNE = "NE"
	ComparisonOperatorGT = "GT"
	ComparisonOperatorLT = "LT"
	ComparisonOperatorGE = "GE"
	ComparisonOperatorLE = "LE"
)

// wrapped struct for condition on Query operation
type QueryCondition struct {
	indexName  string
	conditions map[string]SDK.Condition
}

// Create new QueryCondition struct
func NewQueryCondition() *QueryCondition {
	c := &QueryCondition{}
	c.conditions = make(map[string]SDK.Condition)
	return c
}

// Create new DynamoDB condition for Query operation
func NewCondition(value Any, operator string) SDK.Condition {
	return SDK.Condition{
		AttributeValueList: []*SDK.AttributeValue{createAttributeValue(value)},
		ComparisonOperator: String(operator),
	}
}

// Add a EQUAL condition for Query operation
func (c *QueryCondition) AddEQ(name string, value Any) {
	cond := NewCondition(value, ComparisonOperatorEQ)
	c.Add(name, cond)
}

// Add a NOT EQUAL condition for Query operation
func (c *QueryCondition) AddNE(name string, value Any) {
	cond := NewCondition(value, ComparisonOperatorNE)
	c.Add(name, cond)
}

// Add a GREATER THAN condition for Query operation
func (c *QueryCondition) AddGT(name string, value Any) {
	cond := NewCondition(value, ComparisonOperatorGT)
	c.Add(name, cond)
}

// Add a LESS THAN condition for Query operation
func (c *QueryCondition) AddLT(name string, value Any) {
	cond := NewCondition(value, ComparisonOperatorLT)
	c.Add(name, cond)
}

// Add a GREATER THAN or EQUAL condition for Query operation
func (c *QueryCondition) AddGE(name string, value Any) {
	cond := NewCondition(value, ComparisonOperatorGE)
	c.Add(name, cond)
}

// Add a LESS THAN or EQUAL condition for Query operation
func (c *QueryCondition) AddLE(name string, value Any) {
	cond := NewCondition(value, ComparisonOperatorLE)
	c.Add(name, cond)
}

// Add a condition for Query operation
func (c *QueryCondition) Add(name string, condition SDK.Condition) {
	c.conditions[name] = condition
}

// Set an Indexname for Query operation
func (c *QueryCondition) UseIndex(name string) {
	c.indexName = name
}
