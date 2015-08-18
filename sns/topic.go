// SNS topic

package sns

import (
	SDK "github.com/aws/aws-sdk-go/service/sns"

	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

type SNSTopic struct {
	svc   *AmazonSNS
	name  string
	arn   string
	sound string
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
		TopicArn: String(t.arn),
	})
	if err != nil {
		log.Error("[SNS] error on `Subscribe` operation, topic="+t.arn, err.Error())
		return "", err
	}
	return *resp.SubscriptionArn, nil
}

// Publish notification to the topic
func (t *SNSTopic) Publish(msg string) error {
	return t.svc.Publish(t.arn, msg, nil)
}

// Delete topic
func (t *SNSTopic) Delete() error {
	_, err := t.svc.Client.DeleteTopic(&SDK.DeleteTopicInput{
		TopicArn: String(t.arn),
	})
	return err
}
