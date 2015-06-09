package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testConfig struct{}

func (t *testConfig) GetConfigValue(a, b, c string) string {
	return fmt.Sprintf("test section=%s, key=%s, default=%s", a, b, c)
}
func (t *testConfig) SetValues(v map[string]interface{}) {}
func (t *testConfig) LoadFile(v string) error {
	return nil
}

func TestSetConfig(t *testing.T) {
	assert := assert.New(t)
	assert.Nil(_config)

	SetConfig(&testConfig{})
	assert.Equal(&testConfig{}, _config)
}

func TestGetConfigValue(t *testing.T) {
	assert := assert.New(t)

	SetConfig(&testConfig{})
	v := GetConfigValue("a", "b", "c")
	assert.Equal("test section=a, key=b, default=c", v)
}

func TestParseToString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("foo", ParseToString("foo"))
	assert.Equal("99", ParseToString("99"))
	assert.Equal("99", ParseToString(99))
	assert.Equal("true", ParseToString(true))
	assert.Equal("false", ParseToString(false))

	// be careful about presicion for float
	assert.Equal("99.999000", ParseToString(99.999))
	assert.Equal("-99.999000", ParseToString(-99.999))

	assert.Equal("<nil>", ParseToString(nil))
	assert.Equal("&{}", ParseToString(_config))

	var empty testConfig
	assert.Equal("{}", ParseToString(empty))
}
