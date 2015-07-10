package auth

import (
	"os"
	"testing"

	AWS "github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"

	_ "github.com/evalphobia/aws-sdk-go-wrapper/config/json"
)

func setTestEnv() {
	os.Clearenv()
	os.Setenv("AWS_ACCESS_KEY_ID", "access")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
}

func TestAuth(t *testing.T) {
	setTestEnv()

	a := Auth()
	creds, err := a.Get()
	assert.NotNil(t, creds)
	assert.Nil(t, err)

	auth = nil
}

func TestNewConfig(t *testing.T) {
	assert := assert.New(t)
	os.Clearenv()
	Clear()

	var conf Config
	var awsConf *AWS.Config

	conf = NewConfig("region", "")
	awsConf = conf.Config
	assert.NotNil(awsConf)
	assert.Equal("region", awsConf.Region)
	assert.Empty(awsConf.Endpoint)
	assert.NotNil(awsConf.Credentials)

	creds, err := awsConf.Credentials.Get()
	assert.NotNil(err)

	// with endpoint
	Clear()
	conf = NewConfig("region", "endpoint")
	awsConf = conf.Config
	assert.NotNil(awsConf)
	assert.Equal("region", awsConf.Region)
	assert.Equal("endpoint", awsConf.Endpoint)
	assert.NotNil(awsConf.Credentials)

	creds, err = awsConf.Credentials.Get()
	assert.NotNil(err)

	// from env
	Clear()
	setTestEnv()
	conf = NewConfig("region", "")
	awsConf = conf.Config
	assert.NotNil(awsConf)
	assert.Equal("region", awsConf.Region)
	assert.NotNil(awsConf.Credentials)

	creds, err = awsConf.Credentials.Get()
	assert.Nil(err)
	assert.NotNil(creds)

	// from cache
	os.Clearenv()
	conf = NewConfig("region", "")
	awsConf = conf.Config
	assert.NotNil(awsConf)

	creds, err = awsConf.Credentials.Get()
	assert.Nil(err)
	assert.NotNil(creds)
	auth = nil
}

func TestEnvRegion(t *testing.T) {
	assert := assert.New(t)
	os.Clearenv()

	region := EnvRegion()
	assert.Equal("", region)

	os.Setenv("AWS_REGION", "foobar")
	region = EnvRegion()
	assert.Equal("foobar", region)

	auth = nil
}

func TestNewConfigWithKeys(t *testing.T) {
	assert := assert.New(t)

	var conf Config

	conf = NewConfigWithKeys(Keys{
		AccessKey: "access",
		SecretKey: "secret",
	})
	assert.NotNil(conf.Config)
	assert.NotNil(conf.Config.Credentials)
	assert.Empty(conf.Config.Region)
	assert.Empty(conf.Config.Endpoint)

	conf = NewConfigWithKeys(Keys{
		AccessKey: "access",
		SecretKey: "secret",
		Region:    "region",
		Endpoint:  "endpoint",
	})
	assert.NotNil(conf.Config)
	assert.NotNil(conf.Config.Credentials)
	assert.Equal("region", conf.Config.Region)
	assert.Equal("endpoint", conf.Config.Endpoint)
}

func TestSetDefault(t *testing.T) {
	assert := assert.New(t)

	var conf Config

	conf = NewConfigWithKeys(Keys{
		AccessKey: "access",
		SecretKey: "secret",
	})
	assert.NotNil(conf.Config)
	assert.NotNil(conf.Config.Credentials)
	assert.Empty(conf.Config.Region)
	assert.Empty(conf.Config.Endpoint)

	conf.SetDefault("region", "endpoint")
	assert.Equal("region", conf.Config.Region)
	assert.Equal("endpoint", conf.Config.Endpoint)

	conf = NewConfigWithKeys(Keys{
		AccessKey: "access",
		SecretKey: "secret",
		Region:    "region",
	})
	assert.NotNil(conf.Config)
	assert.NotNil(conf.Config.Credentials)
	assert.Equal("region", conf.Config.Region)
	assert.Empty(conf.Config.Endpoint)

	conf.SetDefault("region", "endpoint")
	assert.Equal("region", conf.Config.Region)
	assert.Empty(conf.Config.Endpoint)
}
