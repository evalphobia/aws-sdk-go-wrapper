package dynamodb

import (
	"reflect"
	"testing"

	SDK "github.com/awslabs/aws-sdk-go/service/dynamodb"
)

func TestNewItem(t *testing.T) {
	item := NewItem()
	if len(item.data) != 0 {
		t.Errorf("DynamoItem.data should be an empty map at initialize time, actual=%v", item.data)
	}
	if len(item.conditions) != 0 {
		t.Errorf("DynamoItem.conditions should be an empty map at initialize time, actual=%v", item.data)
	}
}

func TestAddAttribute(t *testing.T) {
	item := NewItem()
	item.AddAttribute("key", "value")
	if len(item.data) != 1 {
		t.Errorf("DynamoItem.data should have one attribute by adding, actual=%v", item.data)
	}
	_, ok := item.data["key"]
	if !ok {
		t.Errorf("DynamoItem.data could not be added correct data, actual=%v", item.data)
	}

	added, _ := item.data["key"]
	if *added.S != "value" {
		t.Errorf("error on add string value, actual=%+v", added)
	}

	item.AddAttribute("int", 100)
	added, _ = item.data["int"]
	if *added.N != "100" {
		t.Errorf("error on add int value, actual=%+v", added)
	}

	item.AddAttribute("float", 100.99)
	added, _ = item.data["float"]
	if *added.N != "100.99" {
		t.Errorf("error on add int value, actual=%+v", added)
	}

}

func TestAddCondition(t *testing.T) {
	item := NewItem()

	cond := NewExpected()
	item.AddCondition("key", cond)

	if len(item.conditions) != 1 {
		t.Errorf("DynamoItem.conditions should have one condition, actual=%v", item.data)
	}
	_, ok := item.conditions["key"]
	if !ok {
		t.Errorf("DynamoItem.conditions could not be added correct data, actual=%v", item.data)
	}

}

func TestAddConditionExist(t *testing.T) {
	item := NewItem()
	item.AddConditionExist("test_condition")
	cond, _ := item.conditions["test_condition"]
	if *cond.Exists != true {
		t.Errorf("error on AddConditionExist, actual=%+v", cond)
	}
}

func TestAddConditionNotExist(t *testing.T) {
	item := NewItem()
	item.AddConditionNotExist("test_condition")
	cond, _ := item.conditions["test_condition"]
	if *cond.Exists != false {
		t.Errorf("error on AddConditionExist, actual=%+v", cond)
	}
}

func TestAddConditionEQ(t *testing.T) {
	expectedOperator := ComparisonOperatorEQ
	item := NewItem()
	item.AddConditionEQ("int_condition", 99)
	cond, _ := item.conditions["int_condition"]
	if *cond.Value.N != "99" || *cond.ComparisonOperator != expectedOperator {
		t.Errorf("error on AddConditionEQ, actual=%+v", cond)
	}

	item.AddConditionEQ("string_condition", "foo")
	cond, _ = item.conditions["string_condition"]
	if *cond.Value.S != "foo" || *cond.ComparisonOperator != expectedOperator {
		t.Errorf("error on AddConditionEQ, actual=%+v", cond)
	}

	item.AddConditionEQ("bool_condition", true)
	cond, _ = item.conditions["bool_condition"]
	if *cond.Value.BOOL != true || *cond.ComparisonOperator != expectedOperator {
		t.Errorf("error on AddConditionEQ, actual=%+v", cond.Value)
	}
}

func TestAddConditionNE(t *testing.T) {
	expectedOperator := ComparisonOperatorNE
	item := NewItem()
	item.AddConditionNE("int_condition", 99)
	cond, _ := item.conditions["int_condition"]
	if *cond.Value.N != "99" || *cond.ComparisonOperator != expectedOperator {
		t.Errorf("error on AddConditionNE, actual=%+v", cond)
	}
}

func TestAddConditionGT(t *testing.T) {
	expectedOperator := ComparisonOperatorGT
	item := NewItem()
	item.AddConditionGT("int_condition", 99)
	cond, _ := item.conditions["int_condition"]
	if *cond.Value.N != "99" || *cond.ComparisonOperator != expectedOperator {
		t.Errorf("error on AddConditionGT, actual=%+v", cond)
	}
}

func TestAddConditionLT(t *testing.T) {
	expectedOperator := ComparisonOperatorLT
	item := NewItem()
	item.AddConditionLT("int_condition", 99)
	cond, _ := item.conditions["int_condition"]
	if *cond.Value.N != "99" || *cond.ComparisonOperator != expectedOperator {
		t.Errorf("error on AddConditionLT, actual=%+v", cond)
	}
}

func TestAddConditionGE(t *testing.T) {
	expectedOperator := ComparisonOperatorGE
	item := NewItem()
	item.AddConditionGE("int_condition", 99)
	cond, _ := item.conditions["int_condition"]
	if *cond.Value.N != "99" || *cond.ComparisonOperator != expectedOperator {
		t.Errorf("error on AddConditionGE, actual=%+v", cond)
	}
}

func TestAddConditionLE(t *testing.T) {
	expectedOperator := ComparisonOperatorLE
	item := NewItem()
	item.AddConditionLE("int_condition", 99)
	cond, _ := item.conditions["int_condition"]
	if *cond.Value.N != "99" || *cond.ComparisonOperator != expectedOperator {
		t.Errorf("error on AddConditionLE, actual=%+v", cond)
	}
}

func TestNewExpected(t *testing.T) {
	exp := NewExpected()
	if reflect.TypeOf(exp) != reflect.TypeOf(&SDK.ExpectedAttributeValue{}) {
		t.Errorf("error on NewExpected, actual=%+v", exp)
	}
}

func TestNewExpectedCondition(t *testing.T) {
	exp := NewExpectedCondition(99, "foo")
	if *exp.Value.N != "99" || *exp.ComparisonOperator != "foo" {
		t.Errorf("error on NewExpectedCondition, actual=%+v", exp)
	}

}
