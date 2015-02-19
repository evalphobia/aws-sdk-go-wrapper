// SQS client

package sqs

import (
	AWS "github.com/awslabs/aws-sdk-go/aws"
	SQS "github.com/awslabs/aws-sdk-go/gen/sqs"

	"github.com/evalphobia/aws-sdk-go-wrapper/auth"
	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

const (
	sqsConfigSectionName = "sqs"
	defaultRegion        = "us-west-1"
	defaultQueuePrefix   = "dev_"
)

type AmazonSQS struct {
	queues map[string]*Queue
	client *SQS.SQS
}

// Create new AmazonSQS struct
func NewClient() *AmazonSQS {
	s := &AmazonSQS{}
	s.queues = make(map[string]*Queue)
	region := config.GetConfigValue(sqsConfigSectionName, "region", defaultRegion)
	s.client = SQS.New(auth.Auth(), region, nil)
	return s
}

// Get a queue
func (s *AmazonSQS) GetQueue(queue string) (*Queue, error) {
	queueName := GetQueuePrefix() + queue

	// get the queue from cache
	q, ok := s.queues[queueName]
	if ok {
		return q, nil
	}

	// get the queue from server
	url, err := s.client.GetQueueURL(&SQS.GetQueueURLRequest{
		QueueName:              AWS.String(queueName),
		QueueOwnerAWSAccountID: nil,
	})
	if err != nil {
		log.Error("[SQS] error on `GetQueueURL` operation, queue="+queueName, err.Error())
		return nil, err
	}
	q = &Queue{}
	q.name = queueName
	q.url = url.QueueURL
	q.autoDel = false
	q.client = s.client
	s.queues[queueName] = q
	return q, nil
}

// Get the prefix for DynamoDB table
func GetQueuePrefix() string {
	return config.GetConfigValue(sqsConfigSectionName, "prefix", defaultQueuePrefix)
}
