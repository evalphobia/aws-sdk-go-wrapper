package config

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

const defaultRegion = "us-east-1"

// NewSession - Creates a new AWS session given a profile to load from your local creds and a region.
// If region is empty, will get the region from the instance metadata.
func NewSession(profile, region string) (*session.Session, error) {
	if region == "" {
		svc := ec2metadata.New(session.New())
		doc, err := svc.GetInstanceIdentityDocument()
		if err != nil {
			return nil, fmt.Errorf("Failed to stablish AWS session: %w", err)
		}
		if len(doc.AvailabilityZone) > 0 {
			region = doc.AvailabilityZone[:len(doc.AvailabilityZone)-1]
		}
	}
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: aws.String(region)},
		Profile:           profile,
	})
	if err != nil {
		return sess, fmt.Errorf("Failed to stablish AWS session: %w", err)
	}
	return sess, nil
}

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
		Region:      pointers.String(c.getRegion()),
	}

	ep := c.getEndpoint()
	if ep != "" {
		awsConf.Endpoint = &ep
	}

	if c.S3ForcePathStyle {
		awsConf.S3ForcePathStyle = pointers.Bool(true)
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
