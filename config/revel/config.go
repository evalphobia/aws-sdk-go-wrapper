/*
   config library for revel
   use like this
       import "github.com/evalphobia/aws-sdk-go-wrapper/<AWS service name>"
       import _ "github.com/evalphobia/aws-sdk-go-wrapper/config/revel"
*/

package revel

import (
	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	revel_config "github.com/evalphobia/revel-config-loader"
)

const (
	awsConfigFileName    = "aws"
	awsConfigSectionName = "aws"
)

// override loggers in initialize
func init() {
	config.SetConfig(NewRevelConfig())
}

type RevelConfig struct{}

func NewRevelConfig() *RevelConfig {
	return &RevelConfig{}
}

// get prams from config file
func (c *RevelConfig) GetConfigValue(section, key, df string) string {
	return revel_config.GetConfigValueDefault(awsConfigFileName, awsConfigSectionName, section+"."+key, df)
}
