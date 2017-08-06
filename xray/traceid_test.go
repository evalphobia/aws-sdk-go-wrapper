package xray

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextUniqueID(t *testing.T) {
	assert := assert.New(t)

	n := nextUniqueID()
	assert.NotEqual(0, n)

	n2 := nextUniqueID()
	assert.NotEqual(0, n2)
	assert.NotEqual(n, n2)
}

func TestNextID(t *testing.T) {
	assert := assert.New(t)

	n := nextID()
	assert.NotEqual("0", n)

	n2 := nextID()
	assert.NotEqual("0", n2)
	assert.NotEqual(n, n2)
}

func TestNextTraceID(t *testing.T) {
	assert := assert.New(t)

	n := nextTraceID()
	assert.NotEqual("0", n)

	n2 := nextTraceID()
	assert.NotEqual("0", n2)
	assert.NotEqual(n, n2)
}
