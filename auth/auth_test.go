package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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
}

func TestNewConfig(t *testing.T) {
	setTestEnv()

	conf := NewConfig("region")
	assert.NotNil(t, conf)
	assert.Equal(t, "region", conf.Region)
	assert.NotNil(t, conf.Credentials)

	creds, err := conf.Credentials.Get()
	assert.Nil(t, err)
	assert.NotNil(t, creds)
}
