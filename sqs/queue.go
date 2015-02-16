// SQS Queue

package sqs

import (
	AWS "github.com/awslabs/aws-sdk-go/aws"
	SQS "github.com/awslabs/aws-sdk-go/gen/sqs"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"

	"encoding/json"
	"errors"
	"fmt"
)

const (
	defaultMessageIDPrefix = "msg_"
	defaultExpireSecond    = 180
)

// SQS Queue wrapper struct
type Queue struct {
	name         string
	url          AWS.StringValue
	messages     []SQS.SendMessageBatchRequestEntry
	delMessages  []SQS.DeleteMessageBatchRequestEntry
	failedSend   []SQS.BatchResultErrorEntry
	failedDelete []SQS.BatchResultErrorEntry

	autoDel bool
	client  *SQS.SQS
}

// Add message to the send spool
func (q *Queue) AddMessage(message string) {
	m := &SQS.SendMessageBatchRequestEntry{}
	m.MessageBody = AWS.String(message)
	num := fmt.Sprint(len(q.messages) + 1)
	m.ID = AWS.String(defaultMessageIDPrefix + num) // 便宜的に採番して設定
	q.messages = append(q.messages, *m)
}

// Send messages in the send spool
func (q *Queue) Send() error {
	// pack the messages ten each to meet the SQS restriction.
	messages := make(map[int][]SQS.SendMessageBatchRequestEntry)
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
func (q *Queue) send(msg []SQS.SendMessageBatchRequestEntry) error {
	res, err := q.client.SendMessageBatch(&SQS.SendMessageBatchRequest{
		Entries:  msg,
		QueueURL: q.url,
	})
	q.failedSend = append(q.failedSend, res.Failed...)
	return err
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

// Get message from the queue with limit
func (q *Queue) Fetch(num int) (*SQS.ReceiveMessageResult, error) {
	// use long-polling for 1sec when to get multiple messages
	var wait int
	if num > 1 {
		wait = 1
	} else {
		wait = 0
	}

	// receive message from server
	resp, err := q.client.ReceiveMessage(&SQS.ReceiveMessageRequest{
		QueueURL:            q.url,
		WaitTimeSeconds:     AWS.Integer(wait),
		MaxNumberOfMessages: AWS.Integer(num),
		VisibilityTimeout:   AWS.Integer(defaultExpireSecond),
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
func (q *Queue) FetchOne() (*SQS.Message, error) {
	resp, err := q.Fetch(1)
	if err != nil {
		return nil, err
	}
	if len(resp.Messages) == 0 {
		return nil, nil
	}
	return &resp.Messages[0], nil
}

// Get only the body of messages
func (q *Queue) FetchBody(num int) []string {
	var messages []string
	resp, _ := q.Fetch(num)
	if len(resp.Messages) == 0 {
		return messages
	}
	for _, msg := range resp.Messages {
		messages = append(messages, *((&msg).Body))
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
	switch msg.(type) {
	case []SQS.Message:
		for _, m := range msg.([]SQS.Message) {
			q.delMessages = append(q.delMessages, SQS.DeleteMessageBatchRequestEntry{
				ID:            m.MessageID,
				ReceiptHandle: m.ReceiptHandle,
			})
		}
	case SQS.Message:
		q.delMessages = append(q.delMessages, SQS.DeleteMessageBatchRequestEntry{
			ID:            msg.(SQS.Message).MessageID,
			ReceiptHandle: msg.(SQS.Message).ReceiptHandle,
		})
	}
}

// Delete a message from server
func (q *Queue) DeleteMessage(msg *SQS.Message) error {
	err := q.client.DeleteMessage(&SQS.DeleteMessageRequest{
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
	messages := make(map[int][]SQS.DeleteMessageBatchRequestEntry)
	if len(q.delMessages) > 10 {
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
func (q *Queue) delete(msg []SQS.DeleteMessageBatchRequestEntry) error {
	res, err := q.client.DeleteMessageBatch(&SQS.DeleteMessageBatchRequest{
		Entries:  msg,
		QueueURL: q.url,
	})
	if err != nil {
		log.Error("[SQS] error on `DeleteMessageBatch`, queue="+q.name, err.Error())
		q.failedDelete = append(q.failedDelete, res.Failed...)
	}

	return err
}

// Set auto delete flag to this queue
func (q *Queue) AutoDelete(b bool) {
	q.autoDel = b
}
