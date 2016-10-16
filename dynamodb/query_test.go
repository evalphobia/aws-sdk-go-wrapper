package dynamodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupQueryTable(t *testing.T) {
	resetTestTable(t)
	tbl := getTestTable(t)
	for i := 1; i <= 10; i++ {
		putTestTable(tbl, 5, i)
	}
	putTestTable(tbl, 6, 5)
	putTestTable(tbl, 6, 6)
}

func TestQueryEQ(t *testing.T) {
	assert := assert.New(t)
	setupQueryTable(t)
	tbl := getTestTable(t)

	c := tbl.NewConditionList()
	c.AndEQ("id", 5)
	c.AndEQ("time", 3)

	res, err := tbl.Query(c)

	assert.NoError(err)
	assert.EqualValues(res.Count, 1)
	assert.EqualValues(res.ScannedCount, 1)
	assert.NotNil(res)

	m := res.ToSliceMap()
	assert.Equal(m[0]["id"], 5)
	assert.Equal(m[0]["time"], 3)
	assert.Equal(m[0]["lsi_key"], "lsi_value")

	// // only hashkey
	// q2 := tbl.NewConditionList()
	// q2.AndEQ("id", 6)
	// assert.NotNil(q2.table)
	// res, err = q2.Query()

	// assert.NoError(err)
	// assert.EqualValues(2, res.Count)
	// assert.EqualValues(2, res.ScannedCount)
	// assert.Nil(res)
}

func TestQueryLT(t *testing.T) {
	assert := assert.New(t)
	setupQueryTable(t)
	tbl := getTestTable(t)

	c := tbl.NewConditionList()
	c.AndEQ("id", 5)
	c.AndLT("time", 3)

	res, err := tbl.Query(c)

	assert.NoError(err)
	assert.EqualValues(res.Count, 2)
	assert.EqualValues(res.ScannedCount, 2)
	assert.NotNil(res)
}

func TestQueryLE(t *testing.T) {
	assert := assert.New(t)
	setupQueryTable(t)
	tbl := getTestTable(t)

	c := tbl.NewConditionList()
	c.AndEQ("id", 5)
	c.AndLE("time", 3)

	res, err := tbl.Query(c)

	assert.NoError(err)
	assert.EqualValues(res.Count, 3)
	assert.EqualValues(res.ScannedCount, 3)
	assert.NotNil(res)
}

func TestQueryGT(t *testing.T) {
	assert := assert.New(t)
	setupQueryTable(t)
	tbl := getTestTable(t)

	c := tbl.NewConditionList()
	c.AndEQ("id", 5)
	c.AndGT("time", 3)

	res, err := tbl.Query(c)

	assert.NoError(err)
	assert.EqualValues(res.Count, 7)
	assert.EqualValues(res.ScannedCount, 7)
	assert.NotNil(res)
}

func TestQueryGE(t *testing.T) {
	assert := assert.New(t)
	setupQueryTable(t)
	tbl := getTestTable(t)

	c := tbl.NewConditionList()
	c.AndEQ("id", 5)
	c.AndGE("time", 3)

	res, err := tbl.Query(c)

	assert.NoError(err)
	assert.EqualValues(res.Count, 8)
	assert.EqualValues(res.ScannedCount, 8)
	assert.NotNil(res)
}

func TestQueryBETWEEN(t *testing.T) {
	assert := assert.New(t)
	setupQueryTable(t)
	tbl := getTestTable(t)

	c := tbl.NewConditionList()
	c.AndEQ("id", 5)
	c.AndBETWEEN("time", 3, 9)

	res, err := tbl.Query(c)

	assert.NoError(err)
	assert.EqualValues(res.Count, 7)
	assert.EqualValues(res.ScannedCount, 7)
	assert.NotNil(res)
}

func TestQueryWithLimit(t *testing.T) {
	assert := assert.New(t)
	setupQueryTable(t)
	tbl := getTestTable(t)

	c := tbl.NewConditionList()
	c.AndEQ("id", 5)
	c.AndBETWEEN("time", 3, 9)
	c.SetLimit(3)

	res, err := tbl.Query(c)

	assert.NoError(err)
	assert.EqualValues(res.Count, 3)
	assert.EqualValues(res.ScannedCount, 3)
	assert.NotNil(res)
}

func TestQueryCount(t *testing.T) {
	assert := assert.New(t)
	setupQueryTable(t)
	tbl := getTestTable(t)

	c := tbl.NewConditionList()
	c.AndEQ("id", 5)
	c.AndBETWEEN("time", 3, 9)

	res, err := tbl.Count(c)

	assert.NoError(err)
	assert.EqualValues(res.Count, 7)
	assert.EqualValues(res.ScannedCount, 7)
	assert.NotNil(res)
}
