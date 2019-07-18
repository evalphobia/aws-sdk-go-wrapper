package pinpointemail

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
	a.Equal("Pinpoint Email", svc.client.ServiceName)
	a.Equal("https://email.us-east-1.amazonaws.com", svc.client.Endpoint)

	region := "us-west-2"
	svc, err = New(config.Config{
		Region: region,
	})
	a.NoError(err)
	expectedEndpoint := "https://email." + region + ".amazonaws.com"
	a.Equal(expectedEndpoint, svc.client.Endpoint)
}
