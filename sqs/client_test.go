package sqs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	SDK "github.com/aws/aws-sdk-go/service/sqs"
	_ "github.com/evalphobia/aws-sdk-go-wrapper/config/json"
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
	assert := assert.New(t)
	setTestEnv()

	svc := NewClient()
	assert.NotNil(svc.client)
	assert.Equal("sqs", svc.client.ServiceName)
	assert.Equal(defaultEndpoint, svc.client.Endpoint)

	region := "us-west-1"
	os.Setenv("AWS_REGION", region)

	c2 := NewClient()
	endpoint := "https://sqs." + region + ".amazonaws.com"
	assert.Equal(endpoint, c2.client.Endpoint)
}

func TestGetQueue(t *testing.T) {
	assert := assert.New(t)
	setTestEnv()
	createQueue("test")

	svc := NewClient()
	q, err := svc.GetQueue("test")
	assert.Nil(err)
	assert.NotNil(q)

	q, err = svc.GetQueue("non_exist")
	assert.NotNil(err)
	assert.Nil(q)

	// cache
	q, err = svc.GetQueue("test")
	assert.Nil(err)
	assert.NotNil(q)
}

func TestSetQueuePrefix(t *testing.T) {
	assert := assert.New(t)
	svc := NewClient()
	assert.Equal(svc.queuePrefix, defaultQueuePrefix)

	svc2 := NewClient()
	svc2.SetQueuePrefix("test")
	assert.Equal(svc2.queuePrefix, "test")
}

func TestCreateQueueWithName(t *testing.T) {
	assert := assert.New(t)
	svc := NewClient()

	// not exitst
	has, err := svc.IsExistQueue("test2")
	assert.NotNil(err)
	assert.False(has)

	// create
	err = svc.CreateQueueWithName("test2")
	assert.Nil(err)

	// created
	has, err = svc.IsExistQueue("test2")
	assert.Nil(err)
	assert.True(has)
}

func TestIsExistQueue(t *testing.T) {
	assert := assert.New(t)
	createQueue("test")

	svc := NewClient()
	has, err := svc.IsExistQueue("test")
	assert.Nil(err)
	assert.True(has)

	// not exitst
	has, err = svc.IsExistQueue("non-exitst-queue")
	assert.NotNil(err)
	assert.False(has)
}
