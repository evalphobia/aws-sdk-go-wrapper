package dynamodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProvisionedThroughput(t *testing.T) {
	assert := assert.New(t)

	tp := newProvisionedThroughput(80, 600)
	assert.EqualValues(80, *tp.ReadCapacityUnits)
	assert.EqualValues(600, *tp.WriteCapacityUnits)
}
