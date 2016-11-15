package dynamodb

import (
	"fmt"
	"testing"
	"time"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/stretchr/testify/assert"
)

func TestDesign(t *testing.T) {
	assert := assert.New(t)
	tbl := getTestTable(t)

	design, err := tbl.Design()
	assert.NoError(err)
	assert.Equal(tbl.design, design)
}

func TestUpdateThroughput(t *testing.T) {
	assert := assert.New(t)
	tbl := getTestTable(t)

	var r, w int64
	r = 10
	w = 20
	err := tbl.UpdateThroughput(r, w)
	assert.NoError(err)

	time.Sleep(500 * time.Millisecond)
	design, _ := tbl.Design()
	assert.Equal(r, design.readCapacity)
	assert.Equal(w, design.writeCapacity)
}

func TestUpdateReadThroughput(t *testing.T) {
	assert := assert.New(t)
	tbl := getTestTable(t)

	var r int64
	r = 50
	err := tbl.UpdateReadThroughput(r)
	assert.NoError(err)

	time.Sleep(500 * time.Millisecond)
	design, _ := tbl.Design()
	assert.Equal(r, design.readCapacity)
}

func TestUpdateWriteThroughput(t *testing.T) {
	assert := assert.New(t)
	tbl := getTestTable(t)

	var w int64
	w = 30
	err := tbl.UpdateWriteThroughput(w)
	assert.NoError(err)

	time.Sleep(500 * time.Millisecond)
	design, _ := tbl.Design()
	assert.Equal(w, design.writeCapacity)
}

func TestAddItem(t *testing.T) {
	assert := assert.New(t)

	tbl := getTestTable(t)
	item := NewPutItem()
	item.AddAttribute("attr1", 99)
	item.AddConditionEQ("cond1", 5)
	tbl.AddItem(item)
	assert.Len(tbl.putSpool, 1)

	items := tbl.putSpool[0]
	it, ok := items.Item["attr1"]
	assert.True(ok)
	assert.Equal("99", *it.N)

	cond, ok := items.Expected["cond1"]
	assert.True(ok)
	assert.Equal("5", *cond.Value.N)
}

func TestPut(t *testing.T) {
	assert := assert.New(t)

	tbl := getTestTable(t)

	item := NewPutItem()
	item.AddAttribute("id", 100)
	item.AddAttribute("time", 1)
	tbl.AddItem(item)

	err := tbl.Put()
	assert.NoError(err)

	item = NewPutItem()
	item.AddAttribute("id", 100)
	tbl.AddItem(item)
	err = tbl.Put()
	assert.Error(err)
}

func TestGetOne(t *testing.T) {
	assert := assert.New(t)

	tbl := getTestTable(t)
	putTestTable(tbl, 100, 1)

	result, err := tbl.GetOne(100, 1)
	assert.NoError(err)
	assert.Equal(100, result["id"])
	assert.Equal(1, result["time"])
}

func TestScan(t *testing.T) {
	assert := assert.New(t)
	resetTestTable(t)
	tbl := getTestTable(t)
	putTestTable(tbl, 100, 1)

	results, err := tbl.Scan()
	assert.NoError(err)
	assert.Len(results.Items, 1)

	result := results.Items[0]
	assert.Equal("100", *result["id"].N)
	assert.Equal("1", *result["time"].N)
}

func TestScanWithCondition(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		limit                int
		filter               int
		expectedCount        int
		expectedScannedCount int
	}{
		{1, 100, 1, 1},
		{1, 101, 1, 1},
		{1, 102, 1, 1},
		{1, 99, 0, 1},
		{2, 100, 2, 2},
		{3, 100, 3, 3},
		{10, 100, 10, 10},
		{10, 101, 10, 10},
		{10, 102, 10, 10},
		{10, 99, 0, 10},
		{500, 100, 200, 200},
	}

	resetTestTable(t)
	tbl := getTestTable(t)
	for i := 0; i < 100; i++ {
		putTestTable(tbl, 100, i)
		if i%2 == 0 {
			putTestTable(tbl, 101, i)
		} else {
			putTestTable(tbl, 102, i)
		}
	}

	for _, tt := range tests {
		target := fmt.Sprintf("%+v", tt)

		cond := tbl.NewConditionList()
		cond.SetLimit(int64(tt.limit))
		if tt.filter != 0 {
			cond.FilterEQ("id", tt.filter)
		}

		results, err := tbl.ScanWithCondition(cond)
		assert.NoError(err, target)

		switch {
		case tt.filter == 0:
			assert.Len(results.Items, tt.expectedCount, target)
			assert.Equal(tt.expectedCount, int(results.Count), target)
		default:
			assert.True(len(results.Items) <= tt.expectedCount, target)
			assert.True(int(results.Count) <= tt.expectedCount, target)
		}
		assert.EqualValues(tt.expectedScannedCount, results.ScannedCount, target)

		for _, item := range results.ToSliceMap() {
			assert.Equal(tt.filter, item["id"], target)
		}
	}
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	resetTestTable(t)
	tbl := getTestTable(t)
	putTestTable(tbl, 100, 1)

	result, err := tbl.GetOne(100, 1)
	assert.NoError(err)
	assert.Equal(100, result["id"])
	assert.Equal(1, result["time"])

	err = tbl.Delete(result["id"], result["time"])
	assert.NoError(err)

	result, err = tbl.GetOne(100, 1)
	assert.NoError(err)
	assert.Nil(result)
}

func TestForceDeleteAll(t *testing.T) {
	assert := assert.New(t)

	tbl := getTestTable(t)
	putTestTable(tbl, 100, 1)

	item := NewPutItem()
	item.AddAttribute("id", 100)
	item.AddAttribute("time", 2)
	tbl.AddItem(item)
	err := tbl.Put()
	assert.NoError(err)

	result, err := tbl.GetOne(100, 1)
	assert.NoError(err)
	assert.Equal(100, result["id"])
	assert.Equal(1, result["time"])

	err = tbl.ForceDeleteAll()
	assert.NoError(err)
	result, err = tbl.GetOne(100, 1)
	assert.NoError(err)
	assert.Nil(result)

	tbl2 := getTestHashTable(t)
	tbl2.AddItem(item)
	err = tbl2.Put()
	assert.NoError(err)

	err = tbl2.ForceDeleteAll()
	assert.NoError(err)
	result, err = tbl2.GetOne(100)
	assert.NoError(err)
	assert.Nil(result)
}

func TestIsExistPrimaryKeys(t *testing.T) {
	assert := assert.New(t)

	tbl := getTestTable(t)
	putTestTable(tbl, 100, 1)

	item := NewPutItem()
	tbl.AddItem(item)
	err := tbl.validatePutItem(tbl.putSpool[0])
	assert.Error(err)

	item.AddAttribute("id", 100)
	tbl.AddItem(item)
	err = tbl.validatePutItem(tbl.putSpool[1])
	assert.Error(err)

	item.AddAttribute("time", 1)
	tbl.AddItem(item)
	err = tbl.validatePutItem(tbl.putSpool[2])
	assert.NoError(err)
}

func TestTableDesign_CreateTableInput(t *testing.T) {
	assert := assert.New(t)
	const (
		prefix       = "testprefix_"
		tableName    = "table_name"
		tableHashKey = "test_hash_key"
	)
	c := config.Config{
		DefaultPrefix: prefix,
	}

	// input without prefix
	design := NewTableDesignWithHashKeyN(tableName, tableHashKey)
	input := design.CreateTableInput("")
	assert.Equal(tableName, *input.TableName)

	// input with prefix.
	input = design.CreateTableInput(c.DefaultPrefix)
	assert.Equal(prefix+tableName, *input.TableName)
}

func putTestTable(tbl *Table, hValue, rValue interface{}) error {
	item := NewPutItem()
	item.AddAttribute("id", hValue)
	item.AddAttribute("time", rValue)
	item.AddAttribute("lsi_key", "lsi_value")
	tbl.AddItem(item)
	return tbl.Put()
}
