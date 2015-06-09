package revel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	assert := assert.New(t)

	conf := NewConfig()
	assert.NotNil(conf)
}

func TestLoadFile(t *testing.T) {
	t.Skip("todo: implement")
}

func TestSetValues(t *testing.T) {
	t.Skip("todo: implement")
}

func TestGetConfigValue(t *testing.T) {
	t.Skip("todo: write test")
}
