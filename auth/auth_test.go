package auth

import (
	"os"
	"testing"

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
	os.Clearenv()
	conf := NewConfig("region")
	assert.NotNil(t, conf)
	assert.Equal(t, "region", conf.Region)
	assert.NotNil(t, conf.Credentials)

	creds, err := conf.Credentials.Get()
	assert.NotNil(t, err)
	auth = nil

	// from env
	setTestEnv()
	conf = NewConfig("region")
	assert.NotNil(t, conf)
	assert.Equal(t, "region", conf.Region)
	assert.NotNil(t, conf.Credentials)

	creds, err = conf.Credentials.Get()
	assert.Nil(t, err)
	assert.NotNil(t, creds)

	// from cache
	os.Clearenv()
	conf = NewConfig("region")
	assert.NotNil(t, conf)

	creds, err = conf.Credentials.Get()
	assert.Nil(t, err)
	assert.NotNil(t, creds)
	auth = nil
}

func TestEnvRegion(t *testing.T) {
	os.Clearenv()

	region := EnvRegion()
	assert.Equal(t, "", region)

	os.Setenv("AWS_REGION", "foobar")
	region = EnvRegion()
	assert.Equal(t, "foobar", region)

	auth = nil
}
