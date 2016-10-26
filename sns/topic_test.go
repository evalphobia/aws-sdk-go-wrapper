package sns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTopic(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	tp := NewTopic(svc, "arn", "name")
	assert.NotNil(tp)
	assert.Equal("arn", tp.arn)
	assert.Equal("name", tp.name)
	assert.Equal(svc, tp.svc)
}

func TestSubscribe(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	topicName := "fooTopic"
	topic, _ := svc.CreateTopic(topicName)

	e := svc.newApplicationEndpoint("arn")
	res, err := topic.Subscribe(e.arn, e.protocol)
	assert.NoError(err)
	assert.Contains(res, "arn:aws:sns:")
	assert.Contains(res, topicName)
}

func TestTopicPublish(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	if svc.client.Endpoint == defaultEndpoint {
		t.Skip("fakesns does not implement Publish() yet.")
	}

	topicName := "fooTopic"
	topic, _ := svc.CreateTopic(topicName)
	err := topic.Publish("foo")
	assert.NoError(err)
}
