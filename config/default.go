// AWS config libs

package config

import (
	"github.com/evalphobia/aws-sdk-go-wrapper/log"

	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	awsConfigFileName = "aws"
)

var (
	config Config
)

func init() {
	SetDefaultConfig()
}

// get parameter from config file
func SetConfig(conf Config) {
	config = conf
}

func GetConfigValue(section, key, df string) string {
	return config.GetConfigValue(section, key, df)
}

// config interface
type Config interface {
	// params(filename, section, key)
	GetConfigValue(string, string, string) string
}

// struct for default config
type DefaultConfig struct {
	config   map[string]interface{}
	rootPath string
}

func NewDefaultConfig(path string) *DefaultConfig {
	c := &DefaultConfig{}
	c.rootPath = path
	return c
}

func SetDefaultConfig() {
	SetConfig(NewDefaultConfig(""))
}

// get parameter from json file
func (c *DefaultConfig) GetConfigValue(section, key, defaultValue string) string {
	if c.config == nil {
		fileName := c.rootPath + "/" + awsConfigFileName + ".json"
		file, e := ioutil.ReadFile(fileName)
		if e != nil {
			log.Error("[config] DefaultConfig File Load error, file="+fileName, e)
			return defaultValue
		}
		json.Unmarshal(file, &c.config)
	}
	val, ok := c.config[key]
	if ok {
		return ParseToString(val)
	} else {
		return defaultValue
	}
}

// return value
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
