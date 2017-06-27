package mobileanalytics

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
)

const defaultEndpoint = "http://localhost:4568"

func getTestConfig() config.Config {
	return config.Config{
		AccessKey: "access",
		SecretKey: "secret",
		Endpoint:  defaultEndpoint,
	}
}

func getTestClient(t *testing.T) *MobileAnalytics {
	svc, err := New(getTestConfig())
	if err != nil {
		t.Errorf("error on create client; error=%s;", err.Error())
		t.FailNow()
	}
	return svc
}

func TestNewClient(t *testing.T) {
	assert := assert.New(t)

	svc, err := New(getTestConfig())
	assert.NoError(err)
	assert.NotNil(svc.client)
	assert.Equal("mobileanalytics", svc.client.ServiceName)
	assert.Equal(defaultEndpoint, svc.client.Endpoint)

	region := "us-west-1"
	svc, err = New(config.Config{
		Region: region,
	})
	assert.NoError(err)
	expectedEndpoint := "https://mobileanalytics." + region + ".amazonaws.com"
	assert.Equal(expectedEndpoint, svc.client.Endpoint)
}
