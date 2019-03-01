package iam

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
	a.Equal("iam", svc.client.ServiceName)
	a.Equal("https://iam.amazonaws.com", svc.client.Endpoint)

	region := "us-west-2"
	svc, err = New(config.Config{
		Region: region,
	})
	a.NoError(err)
	a.Equal("https://iam.amazonaws.com", svc.client.Endpoint, "IAM endpoint is global")
}
