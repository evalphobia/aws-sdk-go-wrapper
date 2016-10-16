package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvRegion(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		isSuccess bool
		envName   string
		region    string
	}{
		{true, "AWS_REGION", "foo"},
		{true, "AWS_REGION", "bar"},
		{false, "AWS_REGION1", "xxx"},
		{false, "AWS_REGION2", "xxx"},
		{false, "AWS_REGIO", "xxx"},
	}

	defer os.Clearenv()
	for _, tt := range tests {
		target := fmt.Sprintf("%+v", tt)
		os.Clearenv()

		assert.Equal("", EnvRegion(), target)

		os.Setenv(tt.envName, tt.region)
		if !tt.isSuccess {
			assert.Equal("", EnvRegion(), target)
			return
		}

		assert.Equal(tt.region, EnvRegion(), target)
		os.Clearenv()
		assert.Equal("", EnvRegion(), target)
	}
}

func TestEnvEndpoint(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		isSuccess bool
		envName   string
		endpoint  string
	}{
		{true, "AWS_ENDPOINT", "foo"},
		{true, "AWS_ENDPOINT", "bar"},
		{false, "AWS_ENDPOINT1", "xxx"},
		{false, "AWS_ENDPOINT2", "xxx"},
		{false, "AWS_ENDPOIN", "xxx"},
	}

	defer os.Clearenv()
	for _, tt := range tests {
		target := fmt.Sprintf("%+v", tt)
		os.Clearenv()

		assert.Equal("", EnvEndpoint(), target)

		os.Setenv(tt.envName, tt.endpoint)
		if !tt.isSuccess {
			assert.Equal("", EnvEndpoint(), target)
			return
		}

		assert.Equal(tt.endpoint, EnvEndpoint(), target)
		os.Clearenv()
		assert.Equal("", EnvEndpoint(), target)
	}
}
