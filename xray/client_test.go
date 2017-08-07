package xray

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

const defaultEndpoint = "http://localhost:9999"

func getTestConfig() config.Config {
	return config.Config{
		AccessKey: "access",
		SecretKey: "secret",
		Endpoint:  defaultEndpoint,
	}
}

func getTestClient(t *testing.T) *XRay {
	svc, err := New(getTestConfig())
	if err != nil {
		t.Errorf("error on create client; error=%s;", err.Error())
		t.FailNow()
	}
	return svc
}

func TestNew(t *testing.T) {
	assert := assert.New(t)

	svc, err := New(getTestConfig())
	assert.NoError(err)
	assert.NotNil(svc.client)
	assert.Equal("xray", svc.client.ServiceName)
	assert.Equal(defaultEndpoint, svc.client.Endpoint)

	region := "us-west-1"
	svc, err = New(config.Config{
		Region: region,
	})
	assert.NoError(err)
	expectedEndpoint := "https://xray." + region + ".amazonaws.com"
	assert.Equal(expectedEndpoint, svc.client.Endpoint)
}

func TestSetLogger(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)
	assert.Equal(log.DefaultLogger, svc.logger)

	stdLogger := &log.StdLogger{}
	svc.SetLogger(stdLogger)
	assert.Equal(stdLogger, svc.logger)
}
