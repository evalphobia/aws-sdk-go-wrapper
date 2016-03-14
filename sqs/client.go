// SQS client

package sqs

import (
	SDK "github.com/aws/aws-sdk-go/service/sqs"

	"github.com/evalphobia/aws-sdk-go-wrapper/auth"
	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

const (
	sqsConfigSectionName = "sqs"
	defaultRegion        = "us-east-1"
	defaultEndpoint      = "http://localhost:4568"
	defaultQueuePrefix   = "dev_"
)

type AmazonSQS struct {
	queues      map[string]*Queue
	client      *SDK.SQS
	queuePrefix string
}

// Create new AmazonSQS struct
func NewClient() *AmazonSQS {
	svc := &AmazonSQS{}
	svc.queues = make(map[string]*Queue)
	region := config.GetConfigValue(sqsConfigSectionName, "region", auth.EnvRegion())
	endpoint := config.GetConfigValue(sqsConfigSectionName, "endpoint", "")
	conf := auth.NewConfig(region, endpoint)
	conf.SetDefault(defaultRegion, defaultEndpoint)
	svc.client = SDK.New(conf.Config)
	svc.queuePrefix = config.GetConfigValue(sqsConfigSectionName, "prefix", defaultQueuePrefix)
	return svc
}

// Get a queue
func (svc *AmazonSQS) GetQueue(queue string) (*Queue, error) {
	queueName := svc.queuePrefix + queue

	// get the queue from cache
	q, ok := svc.queues[queueName]
	if ok {
		return q, nil
	}

	// get the queue from server
	url, err := svc.client.GetQueueUrl(&SDK.GetQueueUrlInput{
		QueueName:              String(queueName),
		QueueOwnerAWSAccountId: nil,
	})
	if err != nil {
		log.Error("[SQS] error on `GetQueueURL` operation, queue="+queueName, err.Error())
		return nil, err
	}
	q = NewQueue(queueName, url.QueueUrl, svc.client)
	svc.queues[queueName] = q
	return q, nil
}

// Create new SQS Queue
func (svc *AmazonSQS) CreateQueue(in *SDK.CreateQueueInput) error {
	data, err := svc.client.CreateQueue(in)
	if err != nil {
		log.Error("[SQS] Error on `CreateQueue` operation, queue="+*in.QueueName, err)
		return err
	}
	log.Info("[SQS] Complete CreateQueue, queue="+*in.QueueName, *(data.QueueUrl))
	return nil
}

// CreateQueueWithName creates new SQS Queue by the name
func (svc *AmazonSQS) CreateQueueWithName(name string) error {
	return svc.CreateQueue(&SDK.CreateQueueInput{
		QueueName: String(svc.queuePrefix + name),
	})
}

// IsExistQueue check queue
func (svc *AmazonSQS) IsExistQueue(name string) (bool, error) {
	name = svc.queuePrefix + name
	data, err := svc.client.GetQueueUrl(&SDK.GetQueueUrlInput{
		QueueName: String(name),
	})

	switch {
	case err != nil:
		log.Error("[SQS] Error on `GetQueueUrl` operation, queue="+name, err)
		return false, err
	case data == nil:
		return false, nil
	case *data.QueueUrl != "":
		return true, nil
	default:
		return false, nil
	}
}

// SetQueuePrefix set queue prefix
func (svc *AmazonSQS) SetQueuePrefix(queuePrefix string) {
	svc.queuePrefix = queuePrefix
}
