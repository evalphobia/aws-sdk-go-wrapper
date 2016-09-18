package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

const defaultRegion = "us-east-1"

// Config has AWS settings.
type Config struct {
	AccessKey string
	SecretKey string
	Region    string
	Endpoint  string

	// Filename and Profile are used for file credentials
	Filename string
	Profile  string

	// DefaultPrefix is used for service resource prefix
	// e.g.) DynamoDB table, S3 bucket, SQS Queue
	DefaultPrefix string
}

// Session creates AWS session from the Config values.
func (c Config) Session() (*session.Session, error) {
	return session.NewSession(c.AWSConfig())
}

// AWSConfig creates *aws.Config object from the fields.
func (c Config) AWSConfig() *aws.Config {
	cred := c.awsCredentials()
	awsConf := &aws.Config{
		Credentials: cred,
		Region:      stringPtr(c.getRegion()),
	}

	ep := c.getEndpoint()
	if ep != "" {
		awsConf.Endpoint = &ep
	}

	return awsConf
}

func (c Config) awsCredentials() *credentials.Credentials {
	// from env
	cred := credentials.NewEnvCredentials()
	_, err := cred.Get()
	if err == nil {
		return cred
	}

	// from param
	cred = credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, "")
	_, err = cred.Get()
	if err == nil {
		return cred
	}

	// from local file
	return credentials.NewSharedCredentials(c.Filename, c.Profile)
}

func (c Config) getRegion() string {
	if c.Region != "" {
		return c.Region
	}
	reg := EnvRegion()
	if reg != "" {
		return reg
	}
	return defaultRegion
}

func (c Config) getEndpoint() string {
	if c.Endpoint != "" {
		return c.Endpoint
	}
	ep := EnvEndpoint()
	if ep != "" {
		return ep
	}
	return ""
}

func stringPtr(s string) *string {
	return &s
}
