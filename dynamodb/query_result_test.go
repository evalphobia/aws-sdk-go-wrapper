package dynamodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestQueryRsult(t *testing.T) *QueryResult {
	setupQueryTable(t)
	tbl := getTestTable(t)

	c := tbl.NewConditionList()
	c.AndEQ("id", 5)
	c.AndLE("time", 2)
	result, _ := tbl.Query(c)
	return result
}

func TestToSliceMap(t *testing.T) {
	assert := assert.New(t)

	r := getTestQueryRsult(t)

	list := r.ToSliceMap()
	assert.Len(list, 2)

	assert.Equal(5, list[0]["id"])
	assert.Equal(1, list[0]["time"])
	assert.Equal("lsi_value", list[0]["lsi_key"])

	assert.Equal(5, list[1]["id"])
	assert.Equal(2, list[1]["time"])
	assert.Equal("lsi_value", list[1]["lsi_key"])
}

func TestUnmarshal(t *testing.T) {
	assert := assert.New(t)

	type myStruct struct {
		ID       int64
		UnixTime int64  `dynamodb:"time"`
		LSIKey   string `dynamodb:"lsi_key"`
	}

	r := getTestQueryRsult(t)

	var list []*myStruct
	err := r.Unmarshal(&list)
	assert.NoError(err)
	assert.Len(list, 2)

	assert.EqualValues(5, list[0].ID)
	assert.EqualValues(1, list[0].UnixTime)
	assert.Equal("lsi_value", list[0].LSIKey)

	assert.EqualValues(5, list[1].ID)
	assert.EqualValues(2, list[1].UnixTime)
	assert.Equal("lsi_value", list[1].LSIKey)
}
