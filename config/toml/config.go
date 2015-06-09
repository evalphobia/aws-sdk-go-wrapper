package toml

import (
	loader "github.com/evalphobia/go-config-loader"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
)

const (
	configType = "toml"
)

func init() {
	config.SetConfig(NewConfig())
}

// Config is config struct for toml format
type Config struct {
	*loader.Config
}

// NewConfig creates new Config for toml
func NewConfig() *Config {
	return &Config{loader.NewConfig()}
}

// LoadFile loads data from the file of given path
func (c *Config) LoadFile(path string) error {
	return c.LoadConfigs(path, configType)
}

// SetValues adds config values
func (c *Config) SetValues(data map[string]interface{}) {
	c.Update(map[string]interface{}{config.AWSConfigFileName: data})
}

// GetConfigValue gets value from loaded config
func (c *Config) GetConfigValue(section, key, df string) string {
	return c.ValueStringDefault(config.AWSConfigFileName+"."+section+"."+key, df)
}
