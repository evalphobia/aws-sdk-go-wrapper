// SQS client

package sqs

import (
	"strings"
	"sync"

	SDK "github.com/aws/aws-sdk-go/service/sqs"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
	"github.com/evalphobia/aws-sdk-go-wrapper/private/errors"
	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

const (
	serviceName = "SQS"
)

// SQS has SQS client and Queue list.
type SQS struct {
	client *SDK.SQS

	logger log.Logger
	prefix string

	queuesMu sync.RWMutex
	queues   map[string]*Queue
}

// New returns initialized *SQS.
func New(conf config.Config) (*SQS, error) {
	sess, err := conf.Session()
	if err != nil {
		return nil, err
	}

	svc := &SQS{
		client: SDK.New(sess),
		logger: log.DefaultLogger,
		prefix: conf.DefaultPrefix,
		queues: make(map[string]*Queue),
	}
	return svc, nil
}

// SetLogger sets logger.
func (svc *SQS) SetLogger(logger log.Logger) {
	svc.logger = logger
}

// GetQueue gets SQS Queue.
func (svc *SQS) GetQueue(name string) (*Queue, error) {
	queueName := svc.prefix + name

	// get the queue from cache
	svc.queuesMu.RLock()
	q, ok := svc.queues[queueName]
	svc.queuesMu.RUnlock()
	if ok {
		return q, nil
	}

	// get the queue from AWS api.
	url, err := svc.client.GetQueueUrl(&SDK.GetQueueUrlInput{
		QueueName:              pointers.String(queueName),
		QueueOwnerAWSAccountId: nil,
	})
	if err != nil {
		svc.Errorf("error on `GetQueueURL` operation; queue=%s; error=%s;", queueName, err.Error())
		return nil, err
	}

	q = NewQueue(svc, name, *url.QueueUrl)
	svc.queuesMu.Lock()
	svc.queues[queueName] = q
	svc.queuesMu.Unlock()
	return q, nil
}

// CreateQueue creates new SQS Queue.
func (svc *SQS) CreateQueue(in *SDK.CreateQueueInput) error {
	data, err := svc.client.CreateQueue(in)
	if err != nil {
		svc.Errorf("error on `CreateQueue` operation; queue=%s; error=%s;", *in.QueueName, err.Error())
		return err
	}

	svc.Infof("success on `CreateQueue` operation; queue=%s; url=%s;", *in.QueueName, *(data.QueueUrl))
	return nil
}

// CreateQueueWithName creates new SQS Queue by given name
func (svc *SQS) CreateQueueWithName(name string) error {
	queueName := svc.prefix + name
	return svc.CreateQueue(&SDK.CreateQueueInput{
		QueueName: pointers.String(queueName),
	})
}

// IsExistQueue checks if the Queue already exists or not.
func (svc *SQS) IsExistQueue(name string) (bool, error) {
	queueName := svc.prefix + name
	data, err := svc.client.GetQueueUrl(&SDK.GetQueueUrlInput{
		QueueName: pointers.String(queueName),
	})

	switch {
	case isNonExistentQueueError(err):
		return false, nil
	case err != nil:
		svc.Errorf("error on `GetQueueUrl` operation; queue=%s; error=%s", name, err.Error())
		return false, err
	case data == nil:
		return false, nil
	case *data.QueueUrl != "": // queue exists
		return true, nil
	default:
		return false, nil
	}
}

// DeleteQueue detes the SQS Queue.
func (svc *SQS) DeleteQueue(name string) error {
	q, err := svc.GetQueue(name)
	if err != nil {
		return err
	}

	_, err = svc.client.DeleteQueue(&SDK.DeleteQueueInput{
		QueueUrl: q.url,
	})
	if err != nil {
		svc.Errorf("error on `DeleteQueue` operation; queue=%s; error=%s;", name, err.Error())
		return err
	}

	svc.Infof("success on `DeleteQueue` operation; queue=%s; url=%s;", name, *q.url)
	return nil
}

func isNonExistentQueueError(err error) bool {
	const errNonExistentQueue = "NonExistentQueue: "
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), errNonExistentQueue)
}

// Infof logging information.
func (svc *SQS) Infof(format string, v ...interface{}) {
	svc.logger.Infof(serviceName, format, v...)
}

// Errorf logging error information.
func (svc *SQS) Errorf(format string, v ...interface{}) {
	svc.logger.Errorf(serviceName, format, v...)
}

func newErrors() *errors.Errors {
	return errors.NewErrors(serviceName)
}
