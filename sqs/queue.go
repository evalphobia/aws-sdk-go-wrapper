// SQS Queue

package sqs

import (
	"encoding/json"
	"fmt"
	"strconv"

	AWS "github.com/aws/aws-sdk-go/aws"
	SDK "github.com/aws/aws-sdk-go/service/sqs"

	"github.com/evalphobia/aws-sdk-go-wrapper/log"
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
	m.Id = AWS.String(defaultMessageIDPrefix + num) // serial numbering for convenience sake
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

	err := &SQSError{}
	// send message
	for i := 0; i < len(messages); i++ {
		e := q.send(messages[i])
		if e != nil {
			log.Error("[SQS] error on `SendMessageBatch` operation, queue="+q.name, e.Error())
			err.AddMessage(e.Error())
		}
	}

	if err.HasError() {
		return err
	}
	return nil
}

// Send a packed message
func (q *Queue) send(msg []*SDK.SendMessageBatchRequestEntry) error {
	res, err := q.client.SendMessageBatch(&SDK.SendMessageBatchInput{
		Entries:  msg,
		QueueUrl: q.url,
	})
	q.failedSend = append(q.failedSend, res.Failed...)
	return err
}

// Get message from the queue with limit
func (q *Queue) Fetch(num int) ([]*Message, error) {
	// use long-polling for 1sec when to get multiple messages
	var wait int = 0
	if num > 1 {
		wait = 1
	}

	// receive message from server
	resp, err := q.client.ReceiveMessage(&SDK.ReceiveMessageInput{
		QueueUrl:            q.url,
		WaitTimeSeconds:     Long(wait),
		MaxNumberOfMessages: Long(num),
		VisibilityTimeout:   Long(defaultExpireSecond),
	})
	if err != nil {
		log.Error("[SQS] error on `ReceiveMessage` operation, queue="+q.name, err.Error())
	}

	var list []*Message
	if resp == nil || len(resp.Messages) == 0 {
		return list, err
	}

	// delete messages automatically
	if q.autoDel {
		q.AddDeleteList(resp.Messages)
		defer q.DeleteListItems()
	}

	for _, msg := range resp.Messages {
		list = append(list, NewMessage(msg))
	}

	return list, err
}

// Get a single message
func (q *Queue) FetchOne() (*Message, error) {
	msgList, err := q.Fetch(1)
	switch {
	case err != nil:
		return nil, err
	case len(msgList) == 0:
		return nil, nil
	}

	return msgList[0], nil
}

// Get only the body of messages
// ** cannot handle deletion manually as lack of MessageId and ReceiptHandle **
func (q *Queue) FetchBody(num int) []string {
	msgList, err := q.Fetch(num)

	var bodies []string
	switch {
	case err != nil:
		log.Error("[SQS] error on `FetchBody`, queue="+q.name, err.Error())
		return bodies
	case len(msgList) == 0:
		return bodies
	}

	for _, msg := range msgList {
		bodies = append(bodies, msg.Body())
	}

	q.AddDeleteList(msgList)
	if q.autoDel {
		defer q.DeleteListItems()
	}
	return bodies
}

// Get the body of a single message
// ** cannot handle deletion manually as lack of MessageId and ReceiptHandle **
func (q *Queue) FetchBodyOne() string {
	bodies := q.FetchBody(1)
	if len(bodies) == 0 {
		return ""
	}
	return bodies[0]
}

// Add a message to the delete spool
func (q *Queue) AddDeleteList(msg interface{}) {
	switch v := msg.(type) {
	case *SDK.Message:
		q.delMessages = append(q.delMessages, &SDK.DeleteMessageBatchRequestEntry{
			Id:            v.MessageId,
			ReceiptHandle: v.ReceiptHandle,
		})
	case *Message:
		q.delMessages = append(q.delMessages, &SDK.DeleteMessageBatchRequestEntry{
			Id:            v.GetMessageID(),
			ReceiptHandle: v.GetReceiptHandle(),
		})
	case []*SDK.Message:
		for _, m := range v {
			q.AddDeleteList(m)
		}
	case []*Message:
		for _, m := range v {
			q.AddDeleteList(m.message)
		}
	}
}

// Delete a message from server
func (q *Queue) DeleteMessage(msg *Message) error {
	_, err := q.client.DeleteMessage(&SDK.DeleteMessageInput{
		QueueUrl:      q.url,
		ReceiptHandle: msg.GetReceiptHandle(),
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
	if msgCount == 0 {
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

	err := &SQSError{}
	// delete messages
	for i := 0; i < len(messages); i++ {
		e := q.delete(messages[i])
		if e != nil {
			err.AddMessage(e.Error())
		}
	}

	if err.HasError() {
		return err
	}
	return nil
}

// Delete a packed message
func (q *Queue) delete(msg []*SDK.DeleteMessageBatchRequestEntry) error {
	if len(msg) < 1 {
		return nil
	}
	res, err := q.client.DeleteMessageBatch(&SDK.DeleteMessageBatchInput{
		Entries:  msg,
		QueueUrl: q.url,
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
		QueueUrl: q.url,
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
	visible, _ := strconv.Atoi(*m["ApproximateNumberOfMessages"])
	invisible, _ := strconv.Atoi(*m["ApproximateNumberOfMessagesNotVisible"])
	return visible, invisible, nil
}

// Delete all messages on the Queue
func (q *Queue) Purge() error {
	_, err := q.client.PurgeQueue(&SDK.PurgeQueueInput{
		QueueUrl: q.url,
	})
	if err != nil {
		log.Error("[SQS] error on `PurgeQueue`, queue="+q.name, err.Error())
		return err
	}
	return nil
}
