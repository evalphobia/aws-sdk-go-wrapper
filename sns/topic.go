package sns

import (
	SDK "github.com/aws/aws-sdk-go/service/sns"

	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

// Topic is struct for Topic.
type Topic struct {
	svc   *SNS
	name  string
	arn   string
	sound string
}

// NewTopic returns initialized *Topic.
func NewTopic(arn, name string, svc *SNS) *Topic {
	return &Topic{
		arn:   arn,
		name:  name,
		sound: "default",
		svc:   svc,
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
		t.svc.Errorf("error on `Subscribe` operation; name=%s; error=%s;", t.name, err.Error())
		return "", err
	}
	return *resp.SubscriptionArn, nil
}

// Publish sends notification to the topic.
func (t *Topic) Publish(msg string) error {
	return t.svc.Publish(t.arn, msg, nil)
}

// Delete deltes the topic.
func (t *Topic) Delete() error {
	_, err := t.svc.client.DeleteTopic(&SDK.DeleteTopicInput{
		TopicArn: pointers.String(t.arn),
	})
	if err != nil {
		t.svc.Errorf("error on `DeleteTopic` operation; name=%s; error=%s;", t.name, err.Error())
	}
	return err
}
