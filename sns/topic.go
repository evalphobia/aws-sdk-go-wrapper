// SNS client

package sns

import (
	AWS "github.com/awslabs/aws-sdk-go/aws"
	SNS "github.com/awslabs/aws-sdk-go/gen/sns"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

type SNSTopic struct {
	name   string
	arn    string
	sound  string
	client *AmazonSNS
}

// Subscribe
func (t *SNSTopic) Subscribe(endpoint *SNSEndpoint) (string, error) {
	resp, err := t.client.Client.Subscribe(&SNS.SubscribeInput{
		Endpoint: AWS.String(endpoint.arn),
		Protocol: AWS.String(endpoint.protocol),
		TopicARN: AWS.String(t.arn),
	})
	if err != nil {
		log.Error("[SNS] error on `Subscribe` operation, topic="+t.arn, err.Error())
		return "", err
	}
	return *resp.SubscriptionARN, nil
}

// Publish notification to the topic
func (t *SNSTopic) Publish(msg string) error {
	return t.client.Publish(t.arn, msg, nil)
}

// Delete topic
func (t *SNSTopic) Delete() error {
	return t.client.Client.DeleteTopic(&SNS.DeleteTopicInput{AWS.String(t.arn)})
}
