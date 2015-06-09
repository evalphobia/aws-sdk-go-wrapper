package sns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTopic(t *testing.T) {
	setTestEnv()

	svc := NewClient()
	tp := NewTopic("arn", "name", svc)
	assert.NotNil(t, tp)
	assert.Equal(t, "arn", tp.arn)
	assert.Equal(t, "name", tp.name)
	assert.Equal(t, svc, tp.svc)
}

func TestSubscribe(t *testing.T) {
	setTestEnv()

	topicName := "fooTopic"
	svc := NewClient()
	topic, _ := svc.CreateTopic(topicName)

	e := NewEndpoint("arn", "application", svc)
	res, err := topic.Subscribe(e)
	assert.Nil(t, err)
	assert.Contains(t, res, "arn:aws:sns:")
	assert.Contains(t, res, topicName)
}

func TestTopicPublish(t *testing.T) {
	setTestEnv()

	topicName := "fooTopic"
	svc := NewClient()
	topic, _ := svc.CreateTopic(topicName)
	err := topic.Publish("foo")

	t.Skip("fakesns does not implement Publish() yet.")
	_ = err
}
