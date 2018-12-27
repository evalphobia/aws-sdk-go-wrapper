package kms

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
)

func getTestConfig() config.Config {
	return config.Config{
		AccessKey: "access",
		SecretKey: "secret",
	}
}

func TestNew(t *testing.T) {
	a := assert.New(t)

	svc, err := New(getTestConfig())
	a.NoError(err)
	a.NotNil(svc.client)
	a.Equal("kms", svc.client.ServiceName)
	a.Equal("https://kms.us-east-1.amazonaws.com", svc.client.Endpoint)

	region := "us-west-2"
	svc, err = New(config.Config{
		Region: region,
	})
	a.NoError(err)
	expectedEndpoint := "https://kms." + region + ".amazonaws.com"
	a.Equal(expectedEndpoint, svc.client.Endpoint)
}
