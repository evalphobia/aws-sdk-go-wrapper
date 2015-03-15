package dynamodb

import (
	"testing"
)

func TestNewItem(t *testing.T) {
	item := NewItem()
	if len(item.data) != 0 {
		t.Errorf("DynamoItem.data should be an empty map at initialize time, actual=%v", item.data)
	}
	if len(item.conditions) != 0 {
		t.Errorf("DynamoItem.conditions should be an empty map at initialize time, actual=%v", item.data)
	}
	if len(item.counters) != 0 {
		t.Errorf("DynamoItem.counters should be an empty map at initialize time, actual=%v", item.data)
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
}

func TestAddStringValue(t *testing.T) {
	item := NewItem()
	item.AddAttribute("key", "value")
	added, _ := item.data["key"]
	if *added.S != "value" {
		t.Errorf("error on add string value, actual=%+v", added)
	}
}

func TestAddIntValue(t *testing.T) {
	item := NewItem()
	item.AddAttribute("key", 100)
	added, _ := item.data["key"]
	if *added.N != "100" {
		t.Errorf("error on add int value, actual=%+v", added)
	}
}

func TestAddFloatValue(t *testing.T) {
	item := NewItem()
	item.AddAttribute("key", 100.99)
	added, _ := item.data["key"]
	if *added.N != "100.99" {
		t.Errorf("error on add float value, actual=%+v", added)
	}
}
