package sqs

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

func getTestClient(t *testing.T) *SQS {
	svc, err := New(getTestConfig())
	if err != nil {
		t.Errorf("error on create client; error=%s;", err.Error())
		t.FailNow()
	}
	return svc
}

func createQueue(name string) {
	svc, _ := New(getTestConfig())
	ok, _ := svc.IsExistQueue(name)
	if ok {
		return
	}

	svc.CreateQueueWithName(name)
}

func TestNewClient(t *testing.T) {
	assert := assert.New(t)

	svc, err := New(getTestConfig())
	assert.NoError(err)
	assert.NotNil(svc.client)
	assert.Equal("sqs", svc.client.ServiceName)
	assert.Equal(defaultEndpoint, svc.client.Endpoint)

	region := "us-west-1"
	svc, err = New(config.Config{
		Region: region,
	})
	assert.NoError(err)
	expectedEndpoint := "https://sqs." + region + ".amazonaws.com"
	assert.Equal(expectedEndpoint, svc.client.Endpoint)
}

func TestGetQueue(t *testing.T) {
	assert := assert.New(t)
	createQueue("test")
	svc := getTestClient(t)

	q, err := svc.GetQueue("test")
	assert.NoError(err)
	assert.NotNil(q)

	svc.DeleteQueue("non_exist")
	q, err = svc.GetQueue("non_exist")
	assert.Error(err)
	assert.Nil(q)

	// cache
	svc.queues["non_exist"] = svc.queues["test"]
	q, err = svc.GetQueue("non_exist")
	assert.NoError(err)
	assert.NotNil(q)
}

func TestCreateQueueWithName(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	// not exitst
	svc.DeleteQueue("test2")
	has, err := svc.IsExistQueue("test2")
	assert.NoError(err)
	assert.False(has)

	// create
	err = svc.CreateQueueWithName("test2")
	assert.NoError(err)

	// created
	has, err = svc.IsExistQueue("test2")
	assert.NoError(err)
	assert.True(has)
}

func TestIsExistQueue(t *testing.T) {
	assert := assert.New(t)
	createQueue("test")
	svc := getTestClient(t)

	has, err := svc.IsExistQueue("test")
	assert.NoError(err)
	assert.True(has)

	// not exitst
	has, err = svc.IsExistQueue("non-exitst-queue")
	assert.NoError(err)
	assert.False(has)
}

func TestSetPrefix(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	svc.SetPrefix("prefix_")
	ok, _ := svc.IsExistQueue("test")
	if !ok {
		svc.CreateQueueWithName("test")
	}
	// No error
	q, err := svc.GetQueue("test")
	assert.NoError(err)
	assert.NotNil(q)

	// Has error
	svc.SetPrefix("prefix2_")
	q, err = svc.GetQueue("test")
	assert.Error(err)
	assert.Nil(q)

	// No error
	svc.SetPrefix("prefix_")
	q, err = svc.GetQueue("test")
	assert.NoError(err)
	assert.NotNil(q)
}
