package dynamodb

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	_ "github.com/evalphobia/aws-sdk-go-wrapper/config/json"
	"github.com/evalphobia/aws-sdk-go-wrapper/auth"
)

func init() {
	setTestEnv()
}

func setTestEnv() {
	os.Clearenv()
	os.Setenv("AWS_ACCESS_KEY_ID", "access")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
}

func TestNewClient(t *testing.T) {
	assert := assert.New(t)

	c := NewClient()
	assert.NotNil(c)
	assert.NotNil(c.client)
	assert.Len(c.tables, 0)
	assert.Len(c.writeTables, 0)
}

func TestNewClientWithKeys(t *testing.T) {
	assert := assert.New(t)

	c := NewClientWithKeys(auth.Keys{
		AccessKey: "myAccessKey",
		SecretKey: "mySecretKey",
		Region: "ap-northeast-1",
	})
	assert.NotNil(c)
	assert.NotNil(c.client)
	assert.Len(c.tables, 0)
	assert.Len(c.writeTables, 0)
}

func TestCreateTable(t *testing.T) {
	assert := assert.New(t)
	name := "foo_table"

	c := NewClient()
	resetTable(c, name)

	in := NewCreateTableWithHashKeyN(name, "id")
	in.AddRangeKeyN("time")
	in.AddLSIS("lsi-index", "lsi_key")
	in.AddGSINN("gsi-index", "time", "id")

	err := c.CreateTable(in)
	assert.Nil(err)
	
	// duplicate table
	err = c.CreateTable(in)
	assert.NotNil(err)

	// empty request
	empty := CreateTableInput{}
	err = c.CreateTable(&empty)
	assert.NotNil(err)
}

func TestDeleteTable(t *testing.T) {
	assert := assert.New(t)
	name := "foo_table_delete"

	c := NewClient()
	in := getCreateTableInput(name)
	createTable(c, in)

	err := c.DeleteTable(c.TablePrefix + name)
	assert.Nil(err)
	
	// deleted table does not return error
	err = c.DeleteTable(c.TablePrefix + name)
	assert.Nil(err)
}

func TestDeleteTableWithPrefix(t *testing.T) {
	assert := assert.New(t)
	name := "foo_table_delete"

	c := NewClient()
	in := getCreateTableInput(name)
	createTable(c, in)

	err := c.DeleteTableWithPrefix(name)
	assert.Nil(err)

	// deleted table does not return error
	err = c.DeleteTableWithPrefix(name)
	assert.Nil(err)
}

func TestDescribeTable(t *testing.T) {
	assert := assert.New(t)
	name := "foo_table"

	c := NewClient()
	resetTable(c, name)

	in := getCreateTableInput(name)
	createTable(c, in)

	desc, err := c.DescribeTable(c.TablePrefix + name)
	assert.Nil(err)
	assert.NotNil(desc)
	assert.Equal(c.TablePrefix+name, desc.GetTableName())
}

func TestDescribeTableWithPrefix(t *testing.T) {
	assert := assert.New(t)
	name := "foo_table"

	c := NewClient()
	resetTable(c, name)

	in := getCreateTableInput(name)
	createTable(c, in)

	desc, err := c.DescribeTableWithPrefix(name)
	assert.Nil(err)
	assert.NotNil(desc)
	assert.Equal(c.TablePrefix+name, desc.GetTableName())
}

func TestGetTable(t *testing.T) {
	assert := assert.New(t)
	origName := "foo_table"

	c := NewClient()
	resetTable(c, origName)

	in := getCreateTableInput(origName)
	createTable(c, in)

	tbl, err := c.GetTable(origName)
	assert.Nil(err)
	assert.NotNil(tbl)
	assert.NotNil(tbl.table)

	name := c.TablePrefix + origName
	assert.Equal(name, tbl.name)
	assert.Equal(name, *tbl.table.TableName)

	// not exist table
	tbl, err = c.GetTable("non-exist")
	assert.NotNil(err)
	assert.Nil(tbl)

	// get from cache
	tbl, err = c.GetTable(origName)
	assert.Nil(err)
	assert.NotNil(tbl)
}

func TestAddWriteTable(t *testing.T) {
	assert := assert.New(t)
	name := "foo_table"

	c := NewClient()
	c.addWriteTable(name)

	w, ok := c.writeTables[name]
	assert.True(ok)
	assert.Len(c.writeTables, 1)
	assert.True(w)
}

func TestRemoveWriteTable(t *testing.T) {
	assert := assert.New(t)
	name := "foo_table"

	c := NewClient()
	c.removeWriteTable(name)

	w, ok := c.writeTables[name]
	assert.True(ok)
	assert.Len(c.writeTables, 1)
	assert.False(w)
}

func TestListTables(t *testing.T) {
	assert := assert.New(t)
	name := "foo_table"

	c := NewClient()
	resetTable(c, name)

	tables1, err := c.ListTables()
	assert.Nil(err)

	in := getCreateTableInput(name)
	createTable(c, in)

	tables2, err := c.ListTables()
	assert.Nil(err)
	assert.Equal(len(tables1)+1, len(tables2))
}

func TestPutAll(t *testing.T) {
	assert := assert.New(t)
	name := "foo_table"

	c := NewClient()
	resetTable(c, name)
	createTable(c, getCreateTableInput(name))

	tbl, err := c.GetTable(name)
	assert.Nil(err)
	assert.NotNil(tbl)
	assert.Len(tbl.writeItems, 0)

	// add 1 items to tbl
	item := NewItem()
	item.AddAttribute("id", 100)
	item.AddAttribute("time", 1)
	tbl.AddItem(item)

	name2 := "foo_hashtable"
	resetTable(c, name2)
	createTable(c, getCreateHashTableInput(name2))

	tbl2, err := c.GetTable(name2)
	assert.Nil(err)
	assert.NotNil(tbl2)
	assert.Len(tbl2.writeItems, 0)

	// add 3 items to tbl2
	item2a := NewItem()
	item2a.AddAttribute("id", 100)
	item2a.AddAttribute("time", 1)
	item2b := NewItem()
	item2b.AddAttribute("id", 101)
	item2b.AddAttribute("time", 2)
	item2c := NewItem()
	item2c.AddAttribute("id", 102)
	item2c.AddAttribute("time", 3)
	tbl2.AddItem(item2a)
	tbl2.AddItem(item2b)
	tbl2.AddItem(item2c)

	// check before put
	results1, err1 := tbl.Scan()
	results2, err2 := tbl2.Scan()
	assert.Nil(err1)
	assert.Nil(err2)
	assert.Len(results1, 0)
	assert.Len(results2, 0)

	err = c.PutAll()
	assert.Nil(err)

	// check after put
	results1, err1 = tbl.Scan()
	results2, err2 = tbl2.Scan()
	assert.Nil(err1)
	assert.Nil(err2)
	assert.Len(results1, 1)
	assert.Len(results2, 3)
}

func getCreateTableInput(name string) *CreateTableInput {
	in := NewCreateTableWithHashKeyN(name, "id")
	in.AddRangeKeyN("time")
	in.AddLSIS("lsi-index", "lsi_key")
	in.AddGSINN("gsi-index", "time", "id")
	return in
}

func getCreateHashTableInput(name string) *CreateTableInput {
	return NewCreateTableWithHashKeyN(name, "id")
}

func resetTable(c *AmazonDynamoDB, name string) {
	desc, _ := c.DescribeTableWithPrefix(name)
	if desc == nil {
		return
	}
	if desc.IsActive() {
		c.DeleteTableWithPrefix(name)
	}
	time.Sleep(time.Millisecond * 100)
	resetTable(c, name)
}

func createTable(c *AmazonDynamoDB, in *CreateTableInput) {
	desc, _ := c.DescribeTableWithPrefix(in.Name)
	if desc == nil {
		c.CreateTable(in)
		createTable(c, in)
		return
	}
	if desc.IsActive() {
		return
	}
	time.Sleep(time.Millisecond * 100)
	createTable(c, in)
}
