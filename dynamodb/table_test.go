package dynamodb

import (
	// "fmt"
	"testing"
	// "time"

	// SDK "github.com/awslabs/aws-sdk-go/gen/dynamodb"
)

func TestAddItem(t *testing.T) {
	tbl := getTestTable()
	if tbl == nil || len(tbl.writeItems) != 0 {
		t.Errorf("error on GetTable, %v", tbl)
	}

	item := NewItem()
	item.AddAttribute("attr1", 99)
	item.AddConditionEQ("cond1", 5)
	tbl.AddItem(item)
	if len(tbl.writeItems) != 1 {
		t.Errorf("error on AddItem, %v", tbl)
	}
	items := tbl.writeItems[0]

	it, ok := items.Item["attr1"]
	if !ok || *it.N != "99" {
		t.Errorf("error on AddItem, %s", it)
	}

	cond, ok := items.Expected["cond1"]
	if !ok || cond.Value == nil {
		t.Errorf("error on AddItem, %s", cond)
	}
	if cond.Value.N == nil || *cond.Value.N != "5" {
		t.Errorf("error on AddItem, %s", *cond.Value)
	}
}

func TestPut(t *testing.T) {
	tbl := getTestTable()
	if tbl == nil || len(tbl.writeItems) != 0 {
		t.Errorf("error on GetTable, %v", tbl)
	}

	item := NewItem()
	item.AddAttribute("id", 100)
	item.AddAttribute("time", 1)
	tbl.AddItem(item)
	err := tbl.Put()
	if err != nil {
		t.Errorf("error on Put, %s", err.Error())
	}

	item = NewItem()
	item.AddAttribute("id", 100)
	tbl.AddItem(item)
	err = tbl.Put()
	if err == nil {
		t.Errorf("error on Put, %s", err.Error())
	}
}

func TestGetOne(t *testing.T) {
	tbl := getTestTable()
	putTestTable(tbl, 100, 1)

	result, err := tbl.GetOne(100, 1)
	if err != nil {
		t.Errorf("error on GetOne, %s", err.Error())
	}
	if result == nil || result["id"] != 100 || result["time"] != 1 {
		t.Errorf("error on GetOne, %s", result)
	}
}

func TestGetByIndex(t *testing.T) {
	tbl := getTestTable()
	putTestTable(tbl, 100, 1)

	results, err := tbl.GetByIndex("lsi-index", 100, "lsi_value")
	if err != nil {
		t.Errorf("error on GetByIndex, %s", err.Error())
	}
	if len(results) != 1 {
		t.Errorf("error on GetByIndex, %s", results)
	}
	result := results[0]
	if result == nil || result["id"] != 100 || result["time"] != 1 {
		t.Errorf("error on GetByIndex, %s", result)
	}

	results, err = tbl.GetByIndex("gsi-index", 1)
	if err != nil {
		t.Errorf("error on GetByIndex, %s", err.Error())
	}
	if len(results) != 1 {
		t.Errorf("error on GetByIndex, %s", results)
	}
	result = results[0]
	if result == nil || result["id"] != 100 || result["time"] != 1 {
		t.Errorf("error on GetByIndex, %s", result)
	}
}

func TestGet(t *testing.T) {
	// for hashkey+rangekey table
	tbl := getTestTable()
	putTestTable(tbl, 100, 1)
	results, err := tbl.Get(100, 1)
	if err != nil {
		t.Errorf("error on Get, %s", err.Error())
	}
	if len(results) != 1 {
		t.Errorf("error on Get, %s", results)
	}
	result := results[0]
	if result == nil || result["id"] != 100 || result["time"] != 1 {
		t.Errorf("error on Get, %s", result)
	}

	// for hashkey+rangekey table by hashkey condtion
	results, err = tbl.Get(100)
	if err != nil {
		t.Errorf("error on Get, %s", err.Error())
	}
	if len(results) != 1 {
		t.Errorf("error on Get, %s", results)
	}
	result = results[0]
	if result == nil || result["id"] != 100 || result["time"] != 1 {
		t.Errorf("error on Get, %s", result)
	}

	// for hashkey table
	tbl2 := getTestHashTable()
	item := NewItem()
	item.AddAttribute("id", 100)
	tbl2.AddItem(item)
	err = tbl2.Put()
	if err != nil {
		t.Errorf("error on Put, %s", err.Error())
	}

	results, err = tbl2.Get(100)
	if err != nil {
		t.Errorf("error on Get, %s", err.Error())
	}
	if len(results) != 1 {
		t.Errorf("error on Get, %s", results)
	}
	result = results[0]
	if result == nil || result["id"] != 100 {
		t.Errorf("error on Get, %s", result)
	}
}

func TestScan(t *testing.T) {
	tbl := getTestTable()
	putTestTable(tbl, 100, 1)
	results, err := tbl.Scan()
	if err != nil {
		t.Errorf("error on Query, %s", err.Error())
	}
	if len(results) != 1 {
		t.Errorf("error on Query, %s", results)
	}
	result := results[0]
	if result == nil || result["id"] != 100 || result["time"] != 1 {
		t.Errorf("error on Query, %s", result)
	}
}

func TestDelete(t *testing.T) {
	tbl := getTestTable()
	putTestTable(tbl, 100, 1)

	results, err := tbl.Get(100)
	if err != nil {
		t.Errorf("error on Get, %s", err.Error())
	}
	if len(results) != 1 {
		t.Errorf("error on Get, %s", results)
	}

	result := results[0]
	tbl.Delete(result["id"], result["time"])
	results, err = tbl.Get(100)
	if err != nil {
		t.Errorf("error on Delete, %s", err.Error())
	}
	if len(results) != 0 {
		t.Errorf("error on Delete, %s", results)
	}
}

func TestDeleteAll(t *testing.T) {
	tbl := getTestTable()
	putTestTable(tbl, 100, 1)

	item := NewItem()
	item.AddAttribute("id", 100)
	item.AddAttribute("time", 2)
	tbl.AddItem(item)
	err := tbl.Put()
	if err != nil {
		t.Errorf("error on Put, %s", err.Error())
	}

	results, err := tbl.Get(100)
	if err != nil {
		t.Errorf("error on Get, %s", err.Error())
	}
	if len(results) != 2 {
		t.Errorf("error on Get, %s", results)
	}

	tbl.DeleteAll()
	results, err = tbl.Get(100)
	if err != nil {
		t.Errorf("error on DeleteAll, %s", err.Error())
	}
	if len(results) != 0 {
		t.Errorf("error on DeleteAll, %s", results)
	}

	tbl2 := getTestHashTable()
	tbl2.AddItem(item)
	err = tbl2.Put()
	if err != nil {
		t.Errorf("error on Put, %s", err.Error())
	}
	tbl2.DeleteAll()
	results, err = tbl2.Get(100)
	if err != nil {
		t.Errorf("error on DeleteAll, %s", err.Error())
	}
	if len(results) != 0 {
		t.Errorf("error on DeleteAll, %s", results[0])
	}
}

func TestIsExistPrimaryKeys(t *testing.T) {
	tbl := getTestTable()
	putTestTable(tbl, 100, 1)

	item := NewItem()
	tbl.AddItem(item)
	b1 := tbl.isExistPrimaryKeys(tbl.writeItems[0])
	if b1 != false {
		t.Errorf("error on isExistPrimaryKeys, %t", b1)
	}

	item.AddAttribute("id", 100)
	tbl.AddItem(item)
	b2 := tbl.isExistPrimaryKeys(tbl.writeItems[1])
	if b2 != false {
		t.Errorf("error on isExistPrimaryKeys, %t", b2)
	}

	item.AddAttribute("time", 1)
	tbl.AddItem(item)
	b3 := tbl.isExistPrimaryKeys(tbl.writeItems[2])
	if b3 != true {
		t.Errorf("error on isExistPrimaryKeys, %t", b3)
	}
}

func putTestTable(tbl *DynamoTable, hValue, rValue Any) error {
	item := NewItem()
	item.AddAttribute("id", hValue)
	item.AddAttribute("time", rValue)
	item.AddAttribute("lsi_key", "lsi_value")
	tbl.AddItem(item)
	return tbl.Put()
}

func getTestTable() *DynamoTable {
	c := NewClient()
	name := "foo_table"
	in := getCreateTableInput(GetTablePrefix() + name)
	createTable(c, in)
	tbl, _ := c.GetTable(name)
	return tbl
}

func getTestHashTable() *DynamoTable {
	c := NewClient()
	name := "foo_hashtable"
	in := getCreateHashTableInput(GetTablePrefix() + name)
	createTable(c, in)
	tbl, _ := c.GetTable(name)
	return tbl
}
