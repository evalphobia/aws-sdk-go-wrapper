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
	auth *credentials.Credentials = nil
)

// return AWS authorization credentials
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

func NewConfig(region string) *AWS.Config {
	auth := Auth()
	return &AWS.Config{
		Credentials: auth,
		Region:      region,
	}
}

func EnvRegion() string {
	return os.Getenv("AWS_REGION")
}
