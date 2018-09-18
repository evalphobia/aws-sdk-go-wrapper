package sns

import (
	SDK "github.com/aws/aws-sdk-go/service/sns"

	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

// Topic is struct for Topic.
type Topic struct {
	svc            *SNS
	name           string
	nameWithPrefix string
	arn            string
	sound          string
}

// NewTopic returns initialized *Topic.
func NewTopic(svc *SNS, arn, name string) *Topic {
	topicName := svc.prefix + name
	return &Topic{
		svc:            svc,
		arn:            arn,
		name:           name,
		nameWithPrefix: topicName,
		sound:          "default",
	}
}

// Subscribe operates `Subscribe` and returns `SubscriptionArn`.
func (t *Topic) Subscribe(endpointARN, protocol string) (subscriptionARN string, err error) {
	resp, err := t.svc.client.Subscribe(&SDK.SubscribeInput{
		Endpoint: pointers.String(endpointARN),
		Protocol: pointers.String(protocol),
		TopicArn: pointers.String(t.arn),
	})
	if err != nil {
		t.svc.Errorf("error on `Subscribe` operation; name=%s; error=%s;", t.nameWithPrefix, err.Error())
		return "", err
	}
	return *resp.SubscriptionArn, nil
}

// Publish sends notification to the topic.
func (t *Topic) Publish(msg string, isHighPriority bool) error {
	return t.svc.Publish(t.arn, msg, nil, isHighPriority)
}

// Delete deltes the topic.
func (t *Topic) Delete() error {
	_, err := t.svc.client.DeleteTopic(&SDK.DeleteTopicInput{
		TopicArn: pointers.String(t.arn),
	})
	if err != nil {
		t.svc.Errorf("error on `DeleteTopic` operation; name=%s; error=%s;", t.nameWithPrefix, err.Error())
	}
	return err
}
