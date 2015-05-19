package sns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	setTestEnv()

	svc := NewClient()

	app := NewApp("arn", "ios", svc)
	assert.NotNil(t, app)
	assert.Equal(t, "arn", app.arn)
	assert.Equal(t, "ios", app.platform)
	assert.Equal(t, svc, app.svc)
}

func TestCreateEndpoint(t *testing.T) {
	setTestEnv()

	svc := NewClient()
	if svc.Client.Endpoint == defaultEndpoint {
		t.Skip("fakesns does not implement CreatePlatformEndpoint() yet.")
	}

	app := NewApp("arn", "ios", svc)
	ep, err := app.CreateEndpoint("token")
	assert.Nil(t, err)
	assert.NotNil(t, ep)
}
