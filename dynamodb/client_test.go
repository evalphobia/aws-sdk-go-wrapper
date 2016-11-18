package dynamodb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

const (
	defaultEndpoint     = "http://localhost:8000"
	testEmptyBucketName = "test-empty-bucket"
	tablePrefix         = "testprefix_"
)

func getTestConfig() config.Config {
	return config.Config{
		AccessKey: "access",
		SecretKey: "secret",
		Endpoint:  defaultEndpoint,
	}
}

func getTestClient(t *testing.T) *DynamoDB {
	svc, err := New(getTestConfig())
	if err != nil {
		t.Errorf("error on create client; error=%s;", err.Error())
		t.FailNow()
	}
	return svc
}

func TestNew(t *testing.T) {
	assert := assert.New(t)

	svc, err := New(getTestConfig())
	assert.NoError(err)
	assert.NotNil(svc.client)
	assert.Equal("dynamodb", svc.client.ServiceName)
	assert.Equal(defaultEndpoint, svc.client.Endpoint)

	region := "us-west-1"
	svc, err = New(config.Config{
		Region: region,
	})
	assert.NoError(err)
	expectedEndpoint := "https://dynamodb." + region + ".amazonaws.com"
	assert.Equal(expectedEndpoint, svc.client.Endpoint)
}

func TestSetLogger(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)
	assert.Equal(log.DefaultLogger, svc.logger)

	stdLogger := &log.StdLogger{}
	svc.SetLogger(stdLogger)
	assert.Equal(stdLogger, svc.logger)
}

func TestCreateTable(t *testing.T) {
	assert := assert.New(t)
	resetTestTable(t)

	name := "foo_table"
	nameWithPrefix := tablePrefix + name
	svc := getTestClient(t)
	svc.prefix = tablePrefix
	td := NewTableDesignWithHashKeyN(name, "id")
	td.AddRangeKeyN("time")
	td.AddLSIS("lsi-index", "lsi_key")
	td.AddGSINN("gsi-index", "time", "id")

	err := svc.CreateTable(td) // create table which name is "testprefix_foo_table"
	assert.NoError(err, "new table creation should be no error")
	table, err := svc.GetTable(name) // get table which name is "testprefix_foo_table"
	assert.NoError(err, "GetTable should be succeessful when name parameter is \"foo_table\"", name)
	table, err = svc.GetTable(nameWithPrefix) // get table which name is "testprefix_testprefix_foo_table"
	assert.Error(err, "GetTable should fail when name parameter is \"testprefix_foo_table\"", nameWithPrefix)
	assert.Nil(table)

	td.name = name
	err = svc.CreateTable(td) // create table which name is "testprefix_foo_table"
	assert.Error(err, "duplicate creation table should be error", td)

	empty := TableDesign{}
	err = svc.CreateTable(&empty)
	assert.Error(err, "empty table design should be error")
}

func TestForceDeleteTable(t *testing.T) {
	assert := assert.New(t)
	name := "foo_table_delete"

	svc := getTestClient(t)
	td := getTableDesign(name)
	createTable(svc, td)

	err := svc.ForceDeleteTable(name)
	assert.NoError(err)

	// deleted table does not return error
	err = svc.ForceDeleteTable(name)
	assert.Error(err)
}

func TestGetTable(t *testing.T) {
	assert := assert.New(t)
	resetTestTable(t)
	svc := getTestClient(t)

	name := "foo_table"
	in := getTableDesign(name)
	createTable(svc, in)

	tbl, err := svc.GetTable(name)
	assert.NoError(err)
	assert.NotNil(tbl)
	assert.NotNil(tbl.design)

	name = svc.prefix + name
	assert.Equal(name, tbl.name)
	assert.Equal(name, tbl.design.name)

	// not exist table
	tbl, err = svc.GetTable("non-exist")
	assert.Error(err)
	assert.Nil(tbl)

	// get from cache
	tbl, err = svc.GetTable(name)
	assert.NoError(err)
	assert.NotNil(tbl)
}

func TestAddWriteTable(t *testing.T) {
	assert := assert.New(t)
	name := "foo_table"

	svc := getTestClient(t)
	svc.addWriteTable(name)

	_, ok := svc.writeTables[name]
	assert.True(ok)
	assert.Len(svc.writeTables, 1)
}

func TestRemoveWriteTable(t *testing.T) {
	assert := assert.New(t)
	name := "foo_table"

	svc := getTestClient(t)
	svc.writeTables[name] = struct{}{}
	assert.Len(svc.writeTables, 1)

	svc.removeWriteTable(name)

	_, ok := svc.writeTables[name]
	assert.False(ok)
	assert.Len(svc.writeTables, 0)
}

func TestListTables(t *testing.T) {
	assert := assert.New(t)
	resetTestTable(t)

	svc := getTestClient(t)
	tables1, err := svc.ListTables()
	assert.NoError(err)

	// create
	getTestTable(t)

	tables2, err := svc.ListTables()
	assert.NoError(err)
	assert.Equal(len(tables1)+1, len(tables2))
}

func TestPutAll(t *testing.T) {
	assert := assert.New(t)
	resetTestTable(t)
	resetTestHashTable(t)
	getTestTable(t)
	getTestHashTable(t)
	svc := getTestClient(t)

	name := "foo_table"
	tbl, err := svc.GetTable(name)
	assert.NoError(err)
	assert.NotNil(tbl)
	assert.Len(tbl.putSpool, 0)

	// add 1 items to tbl
	item := NewPutItem()
	item.AddAttribute("id", 100)
	item.AddAttribute("time", 1)
	tbl.AddItem(item)

	name2 := "foo_hashtable"
	tbl2, err := svc.GetTable(name2)
	assert.NoError(err)
	assert.NotNil(tbl2)
	assert.Len(tbl2.putSpool, 0)

	// add 3 items to tbl2
	item2a := NewPutItem()
	item2a.AddAttribute("id", 100)
	item2a.AddAttribute("time", 1)
	item2b := NewPutItem()
	item2b.AddAttribute("id", 101)
	item2b.AddAttribute("time", 2)
	item2c := NewPutItem()
	item2c.AddAttribute("id", 102)
	item2c.AddAttribute("time", 3)
	tbl2.AddItem(item2a)
	tbl2.AddItem(item2b)
	tbl2.AddItem(item2c)

	// check before put
	results1, err1 := tbl.Scan()
	results2, err2 := tbl2.Scan()
	assert.NoError(err1)
	assert.NoError(err2)
	assert.Len(results1.Items, 0)
	assert.Len(results2.Items, 0)

	err = svc.PutAll()
	assert.NoError(err)

	// check after put
	results1, err1 = tbl.Scan()
	results2, err2 = tbl2.Scan()
	assert.NoError(err1)
	assert.NoError(err2)
	assert.Len(results1.Items, 1)
	assert.Len(results2.Items, 3)
}

func getTableDesign(name string) *TableDesign {
	in := NewTableDesignWithHashKeyN(name, "id")
	in.AddRangeKeyN("time")
	in.AddLSIS("lsi-index", "lsi_key")
	in.AddGSINN("gsi-index", "time", "id")
	return in
}

func getCreateHashTableInput(name string) *TableDesign {
	return NewTableDesignWithHashKeyN(name, "id")
}

func createTable(svc *DynamoDB, design *TableDesign) {
	tbl, err := svc.GetTable(design.name)
	if err == nil && tbl.design.IsActive() {
		return
	}
	svc.CreateTable(design)
	time.Sleep(time.Millisecond * 100)
	createTable(svc, design)
}

func getTestTable(t *testing.T) *Table {
	svc := getTestClient(t)
	name := "foo_table"
	in := getCreateHashTableInput(name)
	in.AddRangeKeyN("time")
	createTable(svc, in)
	tbl, _ := svc.GetTable(name)
	return tbl
}

func getTestHashTable(t *testing.T) *Table {
	svc := getTestClient(t)
	name := "foo_hashtable"
	in := getCreateHashTableInput(name)
	createTable(svc, in)
	tbl, _ := svc.GetTable(name)
	return tbl
}

func resetTestTable(t *testing.T) {
	const name = "foo_table"
	resetTable(getTestClient(t), name)
	resetTable(getTestClient(t), tablePrefix+name)
}

func resetTestHashTable(t *testing.T) {
	const name = "foo_hashtable"
	resetTable(getTestClient(t), name)
}

func resetTable(svc *DynamoDB, name string) {
	tbl, err := svc.GetTable(name)
	switch {
	case err != nil:
		return
	case tbl == nil:
		return
	case !tbl.design.IsActive():
		return
	}

	err = svc.ForceDeleteTable(name)
	resetTable(svc, name)
}
