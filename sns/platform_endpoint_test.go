package sns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEndpoint(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	ep := svc.newApplicationEndpoint("arn")
	assert.NotNil(ep)
	assert.Equal("arn", ep.arn)
	assert.Equal("application", ep.protocol)
	assert.Equal(svc, ep.svc)
}

func TestEndpointPublish(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	if svc.client.Endpoint == defaultEndpoint {
		t.Skip("fakesns does not implement Publish() yet.")
	}

	ep := svc.newApplicationEndpoint("arn")
	err := ep.Publish("msg", 3)
	assert.NoError(err)
}

func TestGetARN(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	ep := svc.newApplicationEndpoint("arn")
	arn := ep.GetARN()
	assert.Equal(arn, "arn")
}
