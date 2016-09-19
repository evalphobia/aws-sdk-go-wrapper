// SQS Queue

package sqs

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	SDK "github.com/aws/aws-sdk-go/service/sqs"
)

const (
	defaultMessageIDPrefix = "msg_"
	defaultExpireSecond    = 180
)

// Queue is SQS Queue wrapper struct.
type Queue struct {
	service *SQS

	name string
	url  *string

	sendSpoolMu sync.Mutex
	sendSpool   []*SDK.SendMessageBatchRequestEntry
	failedSend  []*SDK.BatchResultErrorEntry

	deleteSpoolMu sync.Mutex
	deleteSpool   []*SDK.DeleteMessageBatchRequestEntry
	failedDelete  []*SDK.BatchResultErrorEntry

	autoDel bool
	expire  int
}

// NewQueue returns initialized *Queue.
func NewQueue(name string, url string, service *SQS) *Queue {
	return &Queue{
		name:    name,
		url:     String(url),
		service: service,
		autoDel: false,
		expire:  defaultExpireSecond,
	}
}

// AutoDelete sets auto delete flag.
func (q *Queue) AutoDelete(b bool) {
	q.autoDel = b
}

// SetExpire sets visibility timeout for message.
func (q *Queue) SetExpire(sec int) {
	q.expire = sec
}

// AddMessage adds message to the send spool.
func (q *Queue) AddMessage(message string) {
	q.sendSpoolMu.Lock()
	defer q.sendSpoolMu.Unlock()

	num := fmt.Sprint(len(q.sendSpool) + 1)
	m := &SDK.SendMessageBatchRequestEntry{
		MessageBody: String(message),
		Id:          String(defaultMessageIDPrefix + num), // serial numbering for convenience sake
	}
	q.sendSpool = append(q.sendSpool, m)
}

// AddMessageJSONMarshal adds message to the send pool with encoding json data.
func (q *Queue) AddMessageJSONMarshal(message interface{}) error {
	msg, err := json.Marshal(message)
	if err != nil {
		q.service.Errorf("error on Queue.AddMessageJSONMarshal `json.Marshal` message=%s; error=%s;", fmt.Sprint(msg), err.Error())
		return err
	}

	q.AddMessage(string(msg))
	return nil
}

// AddMessageMap adds message to the send pool from map data.
func (q *Queue) AddMessageMap(message map[string]interface{}) error {
	return q.AddMessageJSONMarshal(message)
}

// Send sends messages in the send spool
func (q *Queue) Send() error {
	q.sendSpoolMu.Lock()
	defer q.sendSpoolMu.Unlock()

	messages := make(map[int][]*SDK.SendMessageBatchRequestEntry)
	spool := q.sendSpool
	switch {
	case len(spool) > 10:
		for i, msg := range spool {
			v := (i + 1) / 10
			messages[v] = append(messages[v], msg)
		}
	default:
		// pack the messages ten each to follow the SQS restriction.
		messages[0] = append(messages[0], spool...)
	}

	sqsError := &SQSError{}
	// send message
	for i := 0; i < len(messages); i++ {
		err := q.send(messages[i])
		if err != nil {
			q.service.Errorf("error on `SendMessageBatch` operation; queue=%s; error=%s;", q.name, err.Error())
			sqsError.Add(err)
		}
	}
	q.sendSpool = nil

	if sqsError.HasError() {
		return sqsError
	}
	return nil
}

// send operates SendMessageBatchInput ands sends a packed message.
func (q *Queue) send(msg []*SDK.SendMessageBatchRequestEntry) error {
	res, err := q.service.client.SendMessageBatch(&SDK.SendMessageBatchInput{
		Entries:  msg,
		QueueUrl: q.url,
	})
	q.failedSend = append(q.failedSend, res.Failed...)
	return err
}

// Fetch fetches message list from the queue with limit.
func (q *Queue) Fetch(num int) ([]*Message, error) {
	wait := 0

	if num > 1 {
		wait = 1 // use long-polling for 1sec when to get multiple messages
	}

	// receive message from AWS api
	resp, err := q.service.client.ReceiveMessage(&SDK.ReceiveMessageInput{
		QueueUrl:            q.url,
		WaitTimeSeconds:     Long(wait),
		MaxNumberOfMessages: Long(num),
		VisibilityTimeout:   Long(defaultExpireSecond),
	})
	if err != nil {
		q.service.Errorf("error on `ReceiveMessage` operation; queue=%s; error=%s;", q.name, err.Error())
	}

	if resp == nil || len(resp.Messages) == 0 {
		return nil, err
	}

	// delete messages automatically
	if q.autoDel {
		q.AddDeleteList(resp.Messages)
		defer q.DeleteListItems()
	}

	list := make([]*Message, len(resp.Messages))
	for i, msg := range resp.Messages {
		list[i] = NewMessage(msg)
	}
	return list, err
}

// FetchOne fetches a single message.
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

// FetchBody fetches only the body of messages.
// ** cannot handle deletion manually as lack of MessageId and ReceiptHandle **
func (q *Queue) FetchBody(num int) []string {
	msgList, err := q.Fetch(num)
	switch {
	case err != nil:
		return nil
	case len(msgList) == 0:
		return nil
	}

	bodies := make([]string, len(msgList))
	for i, msg := range msgList {
		bodies[i] = msg.Body()
	}

	q.AddDeleteList(msgList)
	if q.autoDel {
		defer q.DeleteListItems()
	}
	return bodies
}

// FetchBodyOne fetches the body of a single message.
// ** cannot handle deletion manually as lack of MessageId and ReceiptHandle **
func (q *Queue) FetchBodyOne() string {
	bodies := q.FetchBody(1)
	if len(bodies) == 0 {
		return ""
	}
	return bodies[0]
}

// AddDeleteList adds a message to the delete spool.
func (q *Queue) AddDeleteList(msg interface{}) {
	switch v := msg.(type) {
	case *SDK.Message:
		q.deleteSpoolMu.Lock()
		defer q.deleteSpoolMu.Unlock()
		q.deleteSpool = append(q.deleteSpool, &SDK.DeleteMessageBatchRequestEntry{
			Id:            v.MessageId,
			ReceiptHandle: v.ReceiptHandle,
		})
	case *Message:
		q.deleteSpoolMu.Lock()
		defer q.deleteSpoolMu.Unlock()
		q.deleteSpool = append(q.deleteSpool, &SDK.DeleteMessageBatchRequestEntry{
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

// DeleteMessage sends the request to AWS api to delete the message.
func (q *Queue) DeleteMessage(msg *Message) error {
	_, err := q.service.client.DeleteMessage(&SDK.DeleteMessageInput{
		QueueUrl:      q.url,
		ReceiptHandle: msg.GetReceiptHandle(),
	})
	if err != nil {
		q.service.Errorf("error on `DeleteMessage`; queue=%s; error=%s;", q.name, err.Error())
	}
	return err
}

// DeleteListItems executes delete operation in the delete spool.
func (q *Queue) DeleteListItems() error {
	q.deleteSpoolMu.Lock()
	defer q.deleteSpoolMu.Unlock()

	// pack the messages ten each to meet the SQS restriction.
	spool := q.deleteSpool
	msgCount := len(q.deleteSpool)
	if msgCount == 0 {
		return nil
	}

	messages := make(map[int][]*SDK.DeleteMessageBatchRequestEntry)
	switch {
	case msgCount > 10:
		for i, msg := range spool {
			v := (i + 1) / 10
			messages[v] = append(messages[v], msg)
		}
	default:
		messages[0] = append(messages[0], q.deleteSpool...)
	}

	// delete messages sequentially
	sqsError := &SQSError{}
	for i := 0; i < len(messages); i++ {
		err := q.delete(messages[i])
		if err != nil {
			sqsError.Add(err)
		}
	}
	q.deleteSpool = nil

	if sqsError.HasError() {
		return sqsError
	}
	return nil
}

// delete operates DeleteMessageBatchInput and deletes a packed message.
func (q *Queue) delete(msg []*SDK.DeleteMessageBatchRequestEntry) error {
	if len(msg) == 0 {
		return nil
	}

	res, err := q.service.client.DeleteMessageBatch(&SDK.DeleteMessageBatchInput{
		Entries:  msg,
		QueueUrl: q.url,
	})
	if err != nil {
		q.service.Errorf("error on `DeleteMessageBatch`; queue=%s; error=%s;", q.name, err.Error())
		q.failedDelete = append(q.failedDelete, res.Failed...)
	}

	return err
}

// CountMessage sends request to AWS api to counts left messages in the Queue.
func (q *Queue) CountMessage() (visible int, invisible int, err error) {
	out, err := q.service.client.GetQueueAttributes(&SDK.GetQueueAttributesInput{
		QueueUrl: q.url,
		AttributeNames: []*string{
			String("ApproximateNumberOfMessages"),
			String("ApproximateNumberOfMessagesNotVisible"),
		},
	})
	if err != nil {
		q.service.Errorf("error on `GetQueueAttributes`; queue=%s; error=%s;", q.name, err.Error())
		return 0, 0, err
	}

	m := out.Attributes
	visible, _ = strconv.Atoi(*m["ApproximateNumberOfMessages"])
	invisible, _ = strconv.Atoi(*m["ApproximateNumberOfMessagesNotVisible"])
	return visible, invisible, nil
}

// Purge deletes all messages in the Queue.
func (q *Queue) Purge() error {
	_, err := q.service.client.PurgeQueue(&SDK.PurgeQueueInput{
		QueueUrl: q.url,
	})
	if err != nil {
		q.service.Errorf("error on `PurgeQueue` operation; queue=%s; error=%s;", q.name, err.Error())
		return err
	}

	q.service.Infof("success on `PurgeQueue` operation; queue=%s;", q.name)
	return nil
}
