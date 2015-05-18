package sns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEndpoint(t *testing.T) {
	setTestEnv()

	svc := NewClient()
	ep := NewEndpoint("arn", "application", svc)
	assert.NotNil(t, ep)
	assert.Equal(t, "arn", ep.arn)
	assert.Equal(t, "application", ep.protocol)
	assert.Equal(t, svc, ep.svc)
}

func TestEndpointPublish(t *testing.T) {
	setTestEnv()

	svc := NewClient()
	ep := NewEndpoint("arn", "application", svc)
	err := ep.Publish("msg", 3)

	t.Skip("fakesns does not implement Publish() yet.")
	_ = err
}

func TestGetARN(t *testing.T) {
	setTestEnv()

	svc := NewClient()
	ep := NewEndpoint("arn", "application", svc)
	arn := ep.GetARN()
	assert.Equal(t, arn, "arn")
}
