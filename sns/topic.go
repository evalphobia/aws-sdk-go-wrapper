// SNS topic

package sns

import (
	SDK "github.com/awslabs/aws-sdk-go/service/sns"

	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

type SNSTopic struct {
	name  string
	arn   string
	sound string
	svc   *AmazonSNS
}

func NewTopic(arn, name string, svc *AmazonSNS) *SNSTopic {
	return &SNSTopic{
		arn:   arn,
		name:  name,
		sound: "default",
		svc:   svc,
	}
}

// Subscribe
func (t *SNSTopic) Subscribe(endpoint *SNSEndpoint) (string, error) {
	resp, err := t.svc.Client.Subscribe(&SDK.SubscribeInput{
		Endpoint: String(endpoint.arn),
		Protocol: String(endpoint.protocol),
		TopicARN: String(t.arn),
	})
	if err != nil {
		log.Error("[SNS] error on `Subscribe` operation, topic="+t.arn, err.Error())
		return "", err
	}
	return *resp.SubscriptionARN, nil
}

// Publish notification to the topic
func (t *SNSTopic) Publish(msg string) error {
	return t.svc.Publish(t.arn, msg, nil)
}

// Delete topic
func (t *SNSTopic) Delete() error {
	_, err := t.svc.Client.DeleteTopic(&SDK.DeleteTopicInput{
		TopicARN: String(t.arn),
	})
	return err
}
