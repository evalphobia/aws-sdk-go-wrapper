// DynamoDB QueryCondition operation/manipuration

package dynamodb

import (
	AWS "github.com/awslabs/aws-sdk-go/aws"
	DynamoDB "github.com/awslabs/aws-sdk-go/gen/dynamodb"
)

// wrapped struct for condition on Query operation
type QueryCondition struct {
	indexName  string
	conditions map[string]DynamoDB.Condition
}

// Create new QueryCondition struct
func NewQueryCondition() *QueryCondition {
	c := &QueryCondition{}
	c.conditions = make(map[string]DynamoDB.Condition)
	return c
}

// Create new DynamoDB condition for Query operation
func NewCondition(value Any, operator string) DynamoDB.Condition {
	return DynamoDB.Condition{
		AttributeValueList: []DynamoDB.AttributeValue{createAttributeValue(value)},
		ComparisonOperator: AWS.String(operator),
	}
}

// Add a EQUAL condition for Query operation
func (c *QueryCondition) AddEQ(name string, value Any) {
	cond := NewCondition(value, DynamoDB.ComparisonOperatorEq)
	c.Add(name, cond)
}

// Add a NOT EQUAL condition for Query operation
func (c *QueryCondition) AddNE(name string, value Any) {
	cond := NewCondition(value, DynamoDB.ComparisonOperatorNe)
	c.Add(name, cond)
}

// Add a GREATER THAN condition for Query operation
func (c *QueryCondition) AddGT(name string, value Any) {
	cond := NewCondition(value, DynamoDB.ComparisonOperatorGt)
	c.Add(name, cond)
}

// Add a LESS THAN condition for Query operation
func (c *QueryCondition) AddLT(name string, value Any) {
	cond := NewCondition(value, DynamoDB.ComparisonOperatorLt)
	c.Add(name, cond)
}

// Add a GREATER THAN or EQUAL condition for Query operation
func (c *QueryCondition) AddGE(name string, value Any) {
	cond := NewCondition(value, DynamoDB.ComparisonOperatorGe)
	c.Add(name, cond)
}

// Add a LESS THAN or EQUAL condition for Query operation
func (c *QueryCondition) AddLE(name string, value Any) {
	cond := NewCondition(value, DynamoDB.ComparisonOperatorLe)
	c.Add(name, cond)
}

// Add a condition for Query operation
func (c *QueryCondition) Add(name string, condition DynamoDB.Condition) {
	c.conditions[name] = condition
}

// Set an Indexname for Query operation
func (c *QueryCondition) UseIndex(name string) {
	c.indexName = name
}
