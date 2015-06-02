package dynamodb

import (
	"os"
	"testing"
	"time"

	SDK "github.com/awslabs/aws-sdk-go/service/dynamodb"
)

func setTestEnv() {
	os.Clearenv()
	os.Setenv("AWS_ACCESS_KEY_ID", "access")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
}

func TestNewClient(t *testing.T) {
	setTestEnv()

	c := NewClient()
	if c == nil || c.client == nil {
		t.Errorf("error on NewClient, actual=%v", c)
	}
	if len(c.tables) != 0 || len(c.writeTables) != 0 {
		t.Errorf("error on NewClient, actual=%v", c)
	}
}

func TestCreateTable(t *testing.T) {
	setTestEnv()

	c := NewClient()
	name := "foo_table"
	in := getCreateTableInput(name)
	resetTable(c, name)

	err := c.CreateTable(&in)
	if err != nil {
		t.Errorf("error on CreateTable, %s", err.Error())
	}
}

func TestDeleteTable(t *testing.T) {
	setTestEnv()

	c := NewClient()
	name := "foo_table_delete"

	in := getCreateTableInput(name)
	createTable(c, in)

	err := c.DeleteTable(name)
	if err != nil {
		t.Errorf("error on DeleteTable, %s", err.Error())
	}
}

func TestDescribeTable(t *testing.T) {
	setTestEnv()

	c := NewClient()
	name := "foo_table"
	resetTable(c, name)

	in := getCreateTableInput(name)
	createTable(c, in)

	desc, err := c.DescribeTable(name)
	if err != nil {
		t.Errorf("error on DescribeTable, %s", err.Error())
	}
	if desc == nil || *desc.TableName != name {
		t.Errorf("error on DescribeTable, %v", desc)
	}
}

func TestGetTable(t *testing.T) {
	setTestEnv()

	c := NewClient()
	origName := "foo_table"
	name := GetTablePrefix() + origName
	resetTable(c, name)

	in := getCreateTableInput(name)
	createTable(c, in)

	tbl, err := c.GetTable(origName)
	if err != nil {
		t.Errorf("error on GetTable, %s", err.Error())
	}
	if tbl == nil || tbl.table == nil {
		t.Errorf("error on GetTable, %v", tbl)
	}
	if tbl.name != name || *tbl.table.TableName != name {
		t.Errorf("error on GetTable, %v", tbl)
	}
}

func TestGetTablePrefix(t *testing.T) {
	setTestEnv()

	pfx := GetTablePrefix()
	if pfx != defaultTablePrefix {
		t.Errorf("error on GetTablePrefix, %s", pfx)
	}
}

func TestAddWriteTable(t *testing.T) {
	setTestEnv()

	c := NewClient()
	name := "foo_table"
	c.addWriteTable(name)
	w, ok := c.writeTables[name]
	if len(c.writeTables) != 1 || !ok {
		t.Errorf("error on addWriteTable, %v", c)
	}
	if w != true {
		t.Errorf("error on addWriteTable, %t", w)
	}
}

func TestRemoveWriteTable(t *testing.T) {
	setTestEnv()

	c := NewClient()
	name := "foo_table"
	c.removeWriteTable(name)
	w, ok := c.writeTables[name]
	if len(c.writeTables) != 1 || !ok {
		t.Errorf("error on removeWriteTable, %v", c)
	}
	if w != false {
		t.Errorf("error on removeWriteTable, %t", w)
	}
}

func TestListTables(t *testing.T) {
	setTestEnv()

	c := NewClient()
	name := "foo_table"
	resetTable(c, name)

	tables1, err := c.ListTables()
	if err != nil {
		t.Errorf("error on ListTables, %s", err.Error())
	}

	in := getCreateTableInput(name)
	createTable(c, in)

	tables2, err := c.ListTables()
	if err != nil {
		t.Errorf("error on ListTables, %s", err.Error())
	}
	if len(tables1)+1 != len(tables2) {
		t.Errorf("error on ListTables, %s, %s", tables1, tables2)
	}
}

func TestPutAll(t *testing.T) {
	setTestEnv()

	c := NewClient()
	pfx := GetTablePrefix()
	name := "foo_table"
	resetTable(c, pfx+name)
	createTable(c, getCreateTableInput(pfx+name))
	tbl, err := c.GetTable(name)
	if err != nil || tbl == nil || len(tbl.writeItems) != 0 {
		t.Errorf("error on GetTable, %v, %v", err, tbl)
	}
	item := NewItem()
	item.AddAttribute("id", 100)
	item.AddAttribute("time", 1)
	tbl.AddItem(item)

	name2 := "foo_hashtable"
	resetTable(c, pfx+name2)
	createTable(c, getCreateHashTableInput(pfx+name2))
	tbl2, err := c.GetTable(name2)
	if err != nil || tbl2 == nil || len(tbl2.writeItems) != 0 {
		t.Errorf("error on GetTable, %v, %v", err, tbl2)
	}
	item2 := NewItem()
	item2.AddAttribute("id", 100)
	item2.AddAttribute("time", 1)
	tbl2.AddItem(item2)

	results1, err1 := tbl.Scan()
	results2, err2 := tbl2.Scan()

	if err1 != nil || err2 != nil || len(results1) > 0 || len(results2) > 0 {
		t.Errorf("error on Scan")
	}

	err = c.PutAll()
	if err != nil {
		t.Errorf("error on PutAll, %s", err.Error())
	}

	results1, err1 = tbl.Scan()
	results2, err2 = tbl2.Scan()
	if err1 != nil || err2 != nil || len(results1) != 1 || len(results2) != 1 {
		t.Errorf("error on PutAll")
	}
}

func getCreateTableInput(name string) SDK.CreateTableInput {
	pKey := NewKeySchema(
		NewHashKeyElement("id"),
		NewRangeKeyElement("time"))

	attrs := NewAttributeDefinitions(
		NewNumberAttribute("id"),
		NewNumberAttribute("time"),
		NewStringAttribute("lsi_key"))

	lsi := NewLSI(
		"lsi-index",
		NewKeySchema(
			NewHashKeyElement("id"),
			NewRangeKeyElement("lsi_key"),
		))

	gsi := NewGSI(
		"gsi-index",
		NewKeySchema(
			NewHashKeyElement("time"),
			NewRangeKeyElement("id")),
		NewProvisionedThroughput(1, 1))

	return SDK.CreateTableInput{
		TableName:              &name,
		KeySchema:              pKey,
		AttributeDefinitions:   attrs,
		LocalSecondaryIndexes:  []*SDK.LocalSecondaryIndex{lsi},
		GlobalSecondaryIndexes: []*SDK.GlobalSecondaryIndex{gsi},
		ProvisionedThroughput:  NewProvisionedThroughput(1, 1),
	}
}

func getCreateHashTableInput(name string) SDK.CreateTableInput {
	pKey := NewKeySchema(
		NewHashKeyElement("id"))

	attrs := NewAttributeDefinitions(
		NewNumberAttribute("id"))
	return SDK.CreateTableInput{
		TableName:             &name,
		KeySchema:             pKey,
		AttributeDefinitions:  attrs,
		ProvisionedThroughput: NewProvisionedThroughput(1, 1),
	}
}

func resetTable(c *AmazonDynamoDB, name string) {
	desc, _ := c.client.DescribeTable(&SDK.DescribeTableInput{
		TableName: &name,
	})
	if desc == nil || desc.Table == nil {
		return
	}
	if *desc.Table.TableStatus == "ACTIVE" {
		c.client.DeleteTable(&SDK.DeleteTableInput{
			TableName: &name,
		})
	}
	time.Sleep(time.Millisecond * 100)
	resetTable(c, name)
}

func createTable(c *AmazonDynamoDB, in SDK.CreateTableInput) {
	desc, _ := c.client.DescribeTable(&SDK.DescribeTableInput{
		TableName: in.TableName,
	})
	if desc == nil || desc.Table == nil {
		c.client.CreateTable(&in)
		createTable(c, in)
		return
	}
	if *desc.Table.TableStatus == "ACTIVE" {
		return
	}
	time.Sleep(time.Millisecond * 100)
	createTable(c, in)
}
