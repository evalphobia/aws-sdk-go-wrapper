package toml

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
	assert := assert.New(t)

	conf := NewConfig()
	assert.Equal("", conf.GetConfigValue("auth", "access_key", ""))
	assert.Equal("", conf.GetConfigValue("empty", "access_key", ""))

	err := conf.LoadFile("./")
	assert.Nil(err)
	assert.NotEmpty(conf.Config)
	assert.Equal("testKey", conf.GetConfigValue("auth", "access_key", ""))
	assert.Equal("", conf.GetConfigValue("empty", "access_key", ""))
}

func TestSetValues(t *testing.T) {
	assert := assert.New(t)

	conf := NewConfig()
	err := conf.LoadFile("./")
	assert.Nil(err)

	assert.Equal("", conf.GetConfigValue("aws", "item1", ""))
	assert.Equal("", conf.GetConfigValue("aws", "item2", ""))

	m := make(map[string]interface{})
	m["item1"] = "foo"
	m["item2"] = 99
	conf.SetValues(map[string]interface{}{"test": m})

	assert.Equal("foo", conf.GetConfigValue("test", "item1", ""))
	assert.Equal("99", conf.GetConfigValue("test", "item2", ""))
}

func TestGetConfigValue(t *testing.T) {
	assert := assert.New(t)

	conf := NewConfig()
	err := conf.LoadFile("./")
	assert.Nil(err)

	assert.NotEmpty(conf.Config)
	assert.Equal("testKey", conf.GetConfigValue("auth", "access_key", ""))
	assert.Equal("dev_", conf.GetConfigValue("sqs", "prefix", ""))

	assert.Equal("", conf.GetConfigValue("empty", "access_key", ""), "empty section")
	assert.Equal("", conf.GetConfigValue("auth", "access_key2", ""), "empty key")
}
