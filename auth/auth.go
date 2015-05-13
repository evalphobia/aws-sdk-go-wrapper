// AWS authorization libs

package auth

import (
	AWS "github.com/awslabs/aws-sdk-go/aws"
	"github.com/evalphobia/aws-sdk-go-wrapper/config"
)

const (
	authConfigSectionName = "auth"
	awsAccessConfigKey    = "access_key"
	awsSecretConfigKey    = "secret_key"
)

var (
	auth *AWS.CredentialsProvider = nil
)

// return AWS authorization credentials
func Auth() *AWS.CredentialsProvider {
	if auth != nil {
		return auth
	}

	// return if environmental params for AWS auth
	env, err := AWS.EnvCreds()
	if err == nil {
		auth = &env
		return auth
	}

	accessKey := config.GetConfigValue(authConfigSectionName, awsAccessConfigKey, "")
	secretKey := config.GetConfigValue(authConfigSectionName, awsSecretConfigKey, "")
	creds := AWS.Creds(accessKey, secretKey, "")
	auth = &creds
	return auth
}

func NewConfig(region string) *AWS.Config {
	auth := Auth()
	return &AWS.Config{
		Credentials: *auth,
		Region:      region,
	}
}

