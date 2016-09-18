package sns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	app := svc.newPlatformApplication("arn", "ios")
	assert.NotNil(app)
	assert.Equal("arn", app.arn)
	assert.Equal("ios", app.platform)
	assert.Equal(svc, app.svc)
}

func TestCreateEndpoint(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	if svc.client.Endpoint == defaultEndpoint {
		t.Skip("fakesns does not implement CreatePlatformEndpoint() yet.")
	}

	app := svc.newPlatformApplication("arn", "ios")
	ep, err := app.CreateEndpoint("token")
	assert.Nil(err)
	assert.NotNil(ep)
}
