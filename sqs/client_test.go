package sqs

import (
	"testing"
	"os"

	SDK "github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
)

func setTestEnv() {
	os.Clearenv()
	os.Setenv("AWS_ACCESS_KEY_ID", "access")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
}

func createQueue(name string) {
	svc := NewClient()
	svc.CreateQueue(&SDK.CreateQueueInput{
		QueueName: String(defaultQueuePrefix + name),
	})
}

func TestNewClient(t *testing.T) {
	setTestEnv()

	svc := NewClient()
	assert.NotNil(t, svc.client)
	assert.Equal(t, "sqs", svc.client.ServiceName)
	assert.Equal(t, defaultEndpoint, svc.client.Endpoint)

	region := "us-west-1"
	os.Setenv("AWS_REGION", region)

	c2 := NewClient()
	endpoint := "https://sqs." + region + ".amazonaws.com"
	assert.Equal(t, endpoint, c2.client.Endpoint)
}

func TestGetQueue(t *testing.T) {
	setTestEnv()
	createQueue("test")

	svc := NewClient()
	q, err := svc.GetQueue("test")
	assert.Nil(t, err)
	assert.NotNil(t, q)

	q, err = svc.GetQueue("non_exist")
	assert.NotNil(t, err)
	assert.Nil(t, q)

	// cache
	q, err = svc.GetQueue("test")
	assert.Nil(t, err)
	assert.NotNil(t, q)
}

func TestGetQueuePrefix(t *testing.T) {
	assert.Equal(t, defaultQueuePrefix, GetQueuePrefix())
}
