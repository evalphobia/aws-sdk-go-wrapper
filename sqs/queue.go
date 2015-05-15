// SQS Queue

package sqs

import (
	AWS "github.com/awslabs/aws-sdk-go/aws"
	SDK "github.com/awslabs/aws-sdk-go/service/sqs"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"

	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

const (
	defaultMessageIDPrefix = "msg_"
	defaultExpireSecond    = 180
)

// SQS Queue wrapper struct
type Queue struct {
	name         string
	url          *string
	messages     []*SDK.SendMessageBatchRequestEntry
	delMessages  []*SDK.DeleteMessageBatchRequestEntry
	failedSend   []*SDK.BatchResultErrorEntry
	failedDelete []*SDK.BatchResultErrorEntry

	autoDel bool
	expire  int
	client  *SDK.SQS
}

func NewQueue(name string, url *string, cli *SDK.SQS) *Queue {
	return &Queue{
		name:    name,
		url:     url,
		client:  cli,
		autoDel: false,
		expire:  defaultExpireSecond,
	}
}

// Set auto delete flag to this queue
func (q *Queue) AutoDelete(b bool) {
	q.autoDel = b
}

// Set visibility timeout for message
func (q *Queue) SetExpire(sec int) {
	q.expire = sec
}

// Add message to the send spool
func (q *Queue) AddMessage(message string) {
	m := &SDK.SendMessageBatchRequestEntry{}
	m.MessageBody = AWS.String(message)
	num := fmt.Sprint(len(q.messages) + 1)
	m.ID = AWS.String(defaultMessageIDPrefix + num) // serial numbering for convenience sake
	q.messages = append(q.messages, m)
}

// Add message spool from map data
func (q *Queue) AddMessageMap(message map[string]interface{}) error {
	msg, err := json.Marshal(message)
	if err != nil {
		log.Error("[SQS] error on `json.Marshal`, msg="+fmt.Sprint(msg), err.Error())
		return err
	}
	q.AddMessage(string(msg))
	return nil
}

// Send messages in the send spool
func (q *Queue) Send() error {
	// pack the messages ten each to meet the SQS restriction.
	messages := make(map[int][]*SDK.SendMessageBatchRequestEntry)
	if len(q.messages) > 10 {
		for i, msg := range q.messages {
			v := (i + 1) / 10
			messages[v] = append(messages[v], msg)
		}
	} else {
		messages[0] = append(messages[0], q.messages...)
	}

	var err error = nil
	errStr := ""
	// send message
	for i := 0; i < len(messages); i++ {
		e := q.send(messages[i])
		if e != nil {
			log.Error("[SQS] error on `SendMessageBatch` operation, queue="+q.name, e.Error())
			errStr = errStr + "," + e.Error()
		}
	}
	if errStr != "" {
		err = errors.New(errStr)
	}
	return err
}

// Send a packed message
func (q *Queue) send(msg []*SDK.SendMessageBatchRequestEntry) error {
	res, err := q.client.SendMessageBatch(&SDK.SendMessageBatchInput{
		Entries:  msg,
		QueueURL: q.url,
	})
	q.failedSend = append(q.failedSend, res.Failed...)
	return err
}

// Get message from the queue with limit
func (q *Queue) Fetch(num int) (*SDK.ReceiveMessageOutput, error) {
	// use long-polling for 1sec when to get multiple messages
	var wait int
	if num > 1 {
		wait = 1
	} else {
		wait = 0
	}

	// receive message from server
	resp, err := q.client.ReceiveMessage(&SDK.ReceiveMessageInput{
		QueueURL:            q.url,
		WaitTimeSeconds:     Long(wait),
		MaxNumberOfMessages: Long(num),
		VisibilityTimeout:   Long(defaultExpireSecond),
	})
	if err != nil {
		log.Error("[SQS] error on `ReceiveMessage` operation, queue="+q.name, err.Error())
	}

	// delete messages automatically
	if q.autoDel && len(resp.Messages) > 0 {
		q.AddDeleteList(resp.Messages)
		defer q.DeleteListItems()
	}

	return resp, err
}

// Get a single message
func (q *Queue) FetchOne() (*SDK.Message, error) {
	resp, err := q.Fetch(1)
	if err != nil {
		return nil, err
	}
	if len(resp.Messages) == 0 {
		return nil, nil
	}
	return resp.Messages[0], nil
}

// Get only the body of messages
func (q *Queue) FetchBody(num int) []string {
	var messages []string
	resp, err := q.Fetch(num)
	if err != nil {
		log.Error("[SQS] error on `FetchBody`, queue="+q.name, err.Error())
		return messages
	}
	if len(resp.Messages) == 0 {
		return messages
	}
	for _, msg := range resp.Messages {
		messages = append(messages, *msg.Body)
	}
	q.AddDeleteList(resp.Messages)
	if q.autoDel {
		defer q.DeleteListItems()
	}
	return messages
}

// Get the body of a single message
func (q *Queue) FetchBodyOne() string {
	messages := q.FetchBody(1)
	if len(messages) == 0 {
		return ""
	}
	return messages[0]
}

// Add a message to the delete spool
func (q *Queue) AddDeleteList(msg interface{}) {
	switch v := msg.(type) {
	case []*SDK.Message:
		for _, m := range v {
			q.delMessages = append(q.delMessages, &SDK.DeleteMessageBatchRequestEntry{
				ID:            m.MessageID,
				ReceiptHandle: m.ReceiptHandle,
			})
		}
	case *SDK.Message:
		q.delMessages = append(q.delMessages, &SDK.DeleteMessageBatchRequestEntry{
			ID:            v.MessageID,
			ReceiptHandle: v.ReceiptHandle,
		})
	}
}

// Delete a message from server
func (q *Queue) DeleteMessage(msg *SDK.Message) error {
	_, err := q.client.DeleteMessage(&SDK.DeleteMessageInput{
		QueueURL:      q.url,
		ReceiptHandle: msg.ReceiptHandle,
	})
	if err != nil {
		log.Error("[SQS] error on `DeleteMessage`, queue="+q.name, err.Error())
	}
	return err
}

// Execute delete operation in the delete spool
func (q *Queue) DeleteListItems() error {
	// pack the messages ten each to meet the SQS restriction.
	msgCount := len(q.delMessages)
	if msgCount < 1 {
		return nil
	}
	messages := make(map[int][]*SDK.DeleteMessageBatchRequestEntry)
	if msgCount > 10 {
		for i, msg := range q.delMessages {
			v := (i + 1) / 10
			messages[v] = append(messages[v], msg)
		}
	} else {
		messages[0] = append(messages[0], q.delMessages...)
	}

	var err error = nil
	errStr := ""
	// delete messages
	for i := 0; i < len(messages); i++ {
		e := q.delete(messages[i])
		if e != nil {
			errStr = errStr + "," + e.Error()
		}
	}
	return err
}

// Delete a packed message
func (q *Queue) delete(msg []*SDK.DeleteMessageBatchRequestEntry) error {
	if len(msg) < 1 {
		return nil
	}
	res, err := q.client.DeleteMessageBatch(&SDK.DeleteMessageBatchInput{
		Entries:  msg,
		QueueURL: q.url,
	})
	if err != nil {
		log.Error("[SQS] error on `DeleteMessageBatch`, queue="+q.name, err.Error())
		q.failedDelete = append(q.failedDelete, res.Failed...)
	}

	return err
}

// Count left messages on the Queue
func (q *Queue) CountMessage() (int, int, error) {
	out, err := q.client.GetQueueAttributes(&SDK.GetQueueAttributesInput{
		QueueURL: q.url,
		AttributeNames: []*string{
			String("ApproximateNumberOfMessages"),
			String("ApproximateNumberOfMessagesNotVisible"),
		},
	})
	if err != nil {
		log.Error("[SQS] error on `GetQueueAttributes`, queue="+q.name, err.Error())
		return 0, 0, err
	}
	m := out.Attributes
	visible, _ := strconv.Atoi(*(*m)["ApproximateNumberOfMessages"])
	invisible, _ := strconv.Atoi(*(*m)["ApproximateNumberOfMessagesNotVisible"])
	return visible, invisible, nil
}

// Delete all messages on the Queue
func (q *Queue) Purge() error {
	_, err := q.client.PurgeQueue(&SDK.PurgeQueueInput{
		QueueURL: q.url,
	})
	if err != nil {
		log.Error("[SQS] error on `PurgeQueue`, queue="+q.name, err.Error())
		return err
	}
	return nil
}
