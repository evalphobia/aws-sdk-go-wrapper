package config

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
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

	// Specific sevice's options
	S3ForcePathStyle bool

	// optional
	CredentialsChainVerboseErrors     bool
	DisableComputeChecksums           bool
	DisableEndpointHostPrefix         bool
	DisableParamValidation            bool
	DisableRestProtocolURICleaning    bool
	DisableSSL                        bool
	EC2MetadataDisableTimeoutOverride bool
	EnableEndpointDiscovery           bool
	EnforceShouldRetryCheck           bool
	LowerCaseHeaderMaps               bool
	S3Disable100Continue              bool
	S3UseAccelerate                   bool
	S3DisableContentMD5Validation     bool
	S3UseARNRegion                    bool
	UseDualStack                      bool

	UseMaxRetries  bool
	MaxRetries     int
	UseConfigCache bool
	muConfigCache  sync.RWMutex
	configCache    *aws.Config
}

// Session creates AWS session from the Config values.
func (c *Config) Session() (*session.Session, error) {
	return session.NewSession(c.AWSConfig())
}

// AWSConfig creates *aws.Config object from the fields.
func (c *Config) AWSConfig() *aws.Config {
	if !c.UseConfigCache {
		return c.awsConfig()
	}

	c.muConfigCache.RLock()
	conf := c.configCache
	c.muConfigCache.RUnlock()

	if conf == nil {
		conf = c.awsConfig()
		c.muConfigCache.Lock()
		c.configCache = conf
		c.muConfigCache.Unlock()
	}
	return conf
}

// awsConfig creates *aws.Config object from the fields.
func (c *Config) awsConfig() *aws.Config {
	cred := c.awsCredentials()
	awsConf := &aws.Config{
		Credentials: cred,
		Region:      pointers.String(c.getRegion()),
	}

	ep := c.getEndpoint()
	if ep != "" {
		awsConf.Endpoint = &ep
	}

	if c.S3ForcePathStyle {
		awsConf.S3ForcePathStyle = pointers.Bool(true)
	}

	if c.CredentialsChainVerboseErrors {
		awsConf.CredentialsChainVerboseErrors = pointers.Bool(true)
	}
	if c.DisableComputeChecksums {
		awsConf.DisableComputeChecksums = pointers.Bool(true)
	}
	if c.DisableEndpointHostPrefix {
		awsConf.DisableEndpointHostPrefix = pointers.Bool(true)
	}
	if c.DisableParamValidation {
		awsConf.DisableParamValidation = pointers.Bool(true)
	}
	if c.DisableRestProtocolURICleaning {
		awsConf.DisableRestProtocolURICleaning = pointers.Bool(true)
	}
	if c.DisableSSL {
		awsConf.DisableSSL = pointers.Bool(true)
	}
	if c.EC2MetadataDisableTimeoutOverride {
		awsConf.EC2MetadataDisableTimeoutOverride = pointers.Bool(true)
	}
	if c.EnableEndpointDiscovery {
		awsConf.EnableEndpointDiscovery = pointers.Bool(true)
	}
	if c.EnforceShouldRetryCheck {
		awsConf.EnforceShouldRetryCheck = pointers.Bool(true)
	}
	if c.LowerCaseHeaderMaps {
		awsConf.LowerCaseHeaderMaps = pointers.Bool(true)
	}
	if c.S3Disable100Continue {
		awsConf.S3Disable100Continue = pointers.Bool(true)
	}
	if c.S3UseAccelerate {
		awsConf.S3UseAccelerate = pointers.Bool(true)
	}
	if c.S3DisableContentMD5Validation {
		awsConf.S3DisableContentMD5Validation = pointers.Bool(true)
	}
	if c.S3UseARNRegion {
		awsConf.S3UseARNRegion = pointers.Bool(true)
	}
	if c.UseDualStack {
		awsConf.UseDualStack = pointers.Bool(true)
	}
	if c.UseMaxRetries && c.MaxRetries >= 0 {
		awsConf.MaxRetries = &c.MaxRetries
	}

	return awsConf
}

func (c *Config) awsCredentials() *credentials.Credentials {
	// from env
	cred := credentials.NewEnvCredentials()
	_, err := cred.Get()
	if err == nil {
		return cred
	}

	// from param
	if c.AccessKey != "" && c.SecretKey != "" {
		return credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, "")
	}

	// from local file
	if c.Filename != "" {
		return credentials.NewSharedCredentials(c.Filename, c.Profile)
	}

	// IAM role
	return nil
}

func (c *Config) getRegion() string {
	if c.Region != "" {
		return c.Region
	}
	reg := EnvRegion()
	if reg != "" {
		return reg
	}
	return defaultRegion
}

func (c *Config) getEndpoint() string {
	if c.Endpoint != "" {
		return c.Endpoint
	}
	ep := EnvEndpoint()
	if ep != "" {
		return ep
	}
	return ""
}
