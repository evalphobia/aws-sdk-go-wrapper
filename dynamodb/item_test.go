package dynamodb

import (
	"testing"

	SDK "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

func TestNewPutItem(t *testing.T) {
	assert := assert.New(t)

	item := NewPutItem()
	assert.Len(item.data, 0)
	assert.Len(item.conditions, 0)
}

func TestAddAttribute(t *testing.T) {
	assert := assert.New(t)

	item := NewPutItem()
	item.AddAttribute("key", "value")
	assert.Len(item.data, 1)

	added, ok := item.data["key"]
	assert.True(ok)
	assert.Equal("value", *added.S)

	item.AddAttribute("int", 100)
	added = item.data["int"]
	assert.Equal("100", *added.N)

	item.AddAttribute("float", 100.99)
	added = item.data["float"]
	assert.Equal("100.99", *added.N)
}

func TestAddCondition(t *testing.T) {
	assert := assert.New(t)

	item := NewPutItem()
	item.addCondition("key", &SDK.ExpectedAttributeValue{})
	assert.Len(item.conditions, 1)

	_, ok := item.conditions["key"]
	assert.True(ok)
}

func TestAddConditionExist(t *testing.T) {
	assert := assert.New(t)

	item := NewPutItem()
	item.AddConditionExist("test_condition")
	cond := item.conditions["test_condition"]
	assert.Equal(true, *cond.Exists)
}

func TestAddConditionNotExist(t *testing.T) {
	assert := assert.New(t)

	item := NewPutItem()
	item.AddConditionNotExist("test_condition")
	cond := item.conditions["test_condition"]
	assert.Equal(false, *cond.Exists)
}

func TestAddConditionEQ(t *testing.T) {
	assert := assert.New(t)

	expectedOperator := ComparisonOperatorEQ
	item := NewPutItem()
	item.AddConditionEQ("int_condition", 99)
	cond := item.conditions["int_condition"]
	assert.Equal("99", *cond.Value.N)
	assert.Equal(expectedOperator, *cond.ComparisonOperator)

	item.AddConditionEQ("string_condition", "foo")
	cond = item.conditions["string_condition"]
	assert.Equal("foo", *cond.Value.S)
	assert.Equal(expectedOperator, *cond.ComparisonOperator)

	item.AddConditionEQ("bool_condition", true)
	cond = item.conditions["bool_condition"]
	assert.Equal(true, *cond.Value.BOOL)
	assert.Equal(expectedOperator, *cond.ComparisonOperator)
}

func TestAddConditionNE(t *testing.T) {
	assert := assert.New(t)

	expectedOperator := ComparisonOperatorNE
	item := NewPutItem()
	item.AddConditionNE("int_condition", 99)
	cond := item.conditions["int_condition"]
	assert.Equal("99", *cond.Value.N)
	assert.Equal(expectedOperator, *cond.ComparisonOperator)
}

func TestAddConditionGT(t *testing.T) {
	assert := assert.New(t)

	expectedOperator := ComparisonOperatorGT
	item := NewPutItem()
	item.AddConditionGT("int_condition", 99)
	cond := item.conditions["int_condition"]
	assert.Equal("99", *cond.Value.N)
	assert.Equal(expectedOperator, *cond.ComparisonOperator)
}

func TestAddConditionLT(t *testing.T) {
	assert := assert.New(t)

	expectedOperator := ComparisonOperatorLT
	item := NewPutItem()
	item.AddConditionLT("int_condition", 99)
	cond := item.conditions["int_condition"]
	assert.Equal("99", *cond.Value.N)
	assert.Equal(expectedOperator, *cond.ComparisonOperator)
}

func TestAddConditionGE(t *testing.T) {
	assert := assert.New(t)

	expectedOperator := ComparisonOperatorGE
	item := NewPutItem()
	item.AddConditionGE("int_condition", 99)
	cond := item.conditions["int_condition"]
	assert.Equal("99", *cond.Value.N)
	assert.Equal(expectedOperator, *cond.ComparisonOperator)
}

func TestAddConditionLE(t *testing.T) {
	assert := assert.New(t)

	expectedOperator := ComparisonOperatorLE
	item := NewPutItem()
	item.AddConditionLE("int_condition", 99)
	cond := item.conditions["int_condition"]
	assert.Equal("99", *cond.Value.N)
	assert.Equal(expectedOperator, *cond.ComparisonOperator)
}

func TestNewExpectedCondition(t *testing.T) {
	assert := assert.New(t)

	exp := NewExpectedCondition(99, "foo")
	assert.IsType(&SDK.ExpectedAttributeValue{}, exp)
	assert.Equal("99", *exp.Value.N)
	assert.Equal("foo", *exp.ComparisonOperator)
}
