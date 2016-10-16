package dynamodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLSI(t *testing.T) {
	assert := assert.New(t)

	keys := NewKeySchema(NewHashKeyElement("foo"))
	lsi := NewLSI("name", keys)
	assert.Equal("name", *lsi.IndexName)
	assert.Equal(ProjectionTypeAll, *lsi.Projection.ProjectionType)

	k := lsi.KeySchema
	assert.Len(k, 1)
	assert.Equal("foo", *k[0].AttributeName)

	lsi2 := NewLSI("name", keys, ProjectionTypeKeysOnly)
	assert.Equal("name", *lsi2.IndexName)
	assert.Equal(ProjectionTypeKeysOnly, *lsi2.Projection.ProjectionType)

	k = lsi2.KeySchema
	assert.Len(k, 1)
	assert.Equal("foo", *k[0].AttributeName)
}

func TestNewGSI(t *testing.T) {
	assert := assert.New(t)

	keys := NewKeySchema(NewHashKeyElement("foo"))
	gsi := NewGSI("name", keys, newProvisionedThroughput(5, 8))
	assert.Equal("name", *gsi.IndexName)
	assert.Equal(ProjectionTypeAll, *gsi.Projection.ProjectionType)

	k := gsi.KeySchema
	assert.Len(k, 1)
	assert.Equal("foo", *k[0].AttributeName)

	tp := gsi.ProvisionedThroughput
	assert.EqualValues(5, *tp.ReadCapacityUnits)
	assert.EqualValues(8, *tp.WriteCapacityUnits)

	gsi2 := NewGSI("name", keys, newProvisionedThroughput(5, 8), ProjectionTypeKeysOnly)
	assert.Equal("name", *gsi2.IndexName)
	assert.Equal(ProjectionTypeKeysOnly, *gsi2.Projection.ProjectionType)

	k = gsi2.KeySchema
	assert.Len(k, 1)
	assert.Equal("foo", *k[0].AttributeName)

	tp = gsi2.ProvisionedThroughput
	assert.EqualValues(5, *tp.ReadCapacityUnits)
	assert.EqualValues(8, *tp.WriteCapacityUnits)
}
