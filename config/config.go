package config

import (
	"fmt"
)

// section key name
const (
	AWSConfigFileName = "aws"
)

var (
	_config Config
)

// Config is interface for config data
type Config interface {
	// params(filename, section, key)
	GetConfigValue(string, string, string) string

	// adds config parameter
	SetValues(map[string]interface{})

	// load config parameter from file
	LoadFile(string) error
}

// SetConfig sets new Config
func SetConfig(conf Config) {
	_config = conf
}

// GetConfigValue gets value from laoded config
func GetConfigValue(section, key, df string) string {
	return _config.GetConfigValue(section, key, df)
}

// ParseToString converts value to string
func ParseToString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case int, int32, int64, uint, uint32, uint64:
		return fmt.Sprintf("%d", t)
	case float32, float64:
		return fmt.Sprintf("%f", t)
	case bool:
		return fmt.Sprintf("%t", t)
	case nil:
		return "<nil>"
	default:
		return fmt.Sprintf("%+v", t)
	}
}
