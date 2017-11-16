package sns

import (
	"errors"
	"fmt"
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

func TestParseARNFromError(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		expected    bool
		expectedARN string
		text        string
	}{
		{true, "arn:aws:sns:", "Endpoint arn:aws:sns: already exists"},
		{true, "arn:aws:sns:123456", "Endpoint arn:aws:sns:123456 already exists"},
		{true, "arn:aws:sns:123456", "anbcdEndpoint arn:aws:sns:123456 already exists...."},
		{true, "arn:aws:sns:@#a-b-45_ads", "Endpoint arn:aws:sns:@#a-b-45_ads already exists"},
		{false, "", "ndpoint arn:aws:sns: already exists"},
		{false, "", "Endpoint arn:aws:sns: already exist"},
		{false, "", "Endpoint arn:aws:sns already exists"},
		{false, "", "Endpoint rn:aws:sns: already exists"},
	}

	for _, tt := range tests {
		target := fmt.Sprintf("%+v", tt)

		err := errors.New(tt.text)
		arn, ok := ParseARNFromError(err)
		a.Equal(tt.expected, ok, target)
		a.Equal(tt.expectedARN, arn, target)
	}

	arn, ok := ParseARNFromError(nil)
	a.Equal(false, ok, "When error=nil")
	a.Equal("", arn, "When error=nil")
}
