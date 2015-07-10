// AWS authorization libs

package auth

import (
	"os"

	AWS "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/evalphobia/aws-sdk-go-wrapper/config"
)

const (
	authConfigSectionName = "auth"
	awsAccessConfigKey    = "access_key"
	awsSecretConfigKey    = "secret_key"
)

var (
	auth *credentials.Credentials
)

// Auth return AWS authorization credentials
func Auth() *credentials.Credentials {
	if auth != nil {
		return auth
	}

	// return if environmental params for AWS auth
	e := credentials.NewEnvCredentials()
	_, err := e.Get()
	if err == nil {
		auth = e
		return auth
	}

	accessKey := config.GetConfigValue(authConfigSectionName, awsAccessConfigKey, "")
	secretKey := config.GetConfigValue(authConfigSectionName, awsSecretConfigKey, "")
	auth = credentials.NewStaticCredentials(accessKey, secretKey, "")
	return auth
}

// NewConfig returns initialized Config
func NewConfig(region, endpoint string) Config {
	auth := Auth()
	awsConf := &AWS.Config{
		Credentials: auth,
		Region:      region,
		Endpoint:    endpoint,
	}
	return Config{awsConf}
}

// Clear deletes cache for auth
func Clear() {
	auth = nil
}

// EnvRegion get region from env params
func EnvRegion() string {
	return os.Getenv("AWS_REGION")
}

// Keys used for manual initialization of config on NewConfigWithKeys(Key)
type Keys struct {
	AccessKey string
	SecretKey string
	Region    string
	Endpoint  string
}

// NewConfigWithKeys returns initialized Config with given parameters
func NewConfigWithKeys(k Keys) Config {
	auth := credentials.NewStaticCredentials(k.AccessKey, k.SecretKey, "")
	awsConf := &AWS.Config{
		Credentials: auth,
		Region:      k.Region,
		Endpoint:    k.Endpoint,
	}
	return Config{awsConf}
}

// Config is wrapper struct of AWS.Config
type Config struct {
	*AWS.Config
}

// SetDefault fills parameter of region and endpoint when empty
func (c Config) SetDefault(region, endpoint string) {
	awsConf := c.Config
	if awsConf.Region == "" && awsConf.Endpoint == "" {
		awsConf.Region = region
		awsConf.Endpoint = endpoint
	}
}
