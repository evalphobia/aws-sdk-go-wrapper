// SQS Queue

package sqs

import (
	"encoding/json"
	"fmt"
	"sync"

	SDK "github.com/aws/aws-sdk-go/service/sqs"

	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

const (
	defaultMessageIDPrefix = "msg_"
	defaultExpireSecond    = 180
	defaultWaitTimeSeconds = 0
)

// Queue is SQS Queue wrapper struct.
type Queue struct {
	service *SQS

	name           string
	nameWithPrefix string
	url            *string

	sendSpoolMu sync.Mutex
	sendSpool   []*SDK.SendMessageBatchRequestEntry
	failedSend  []*SDK.BatchResultErrorEntry

	deleteSpoolMu sync.Mutex
	deleteSpool   []*SDK.DeleteMessageBatchRequestEntry
	failedDelete  []*SDK.BatchResultErrorEntry

	autoDel         bool
	expire          int
	waitTimeSeconds int
}

// NewQueue returns initialized *Queue.
func NewQueue(svc *SQS, name string, url string) *Queue {
	queueName := svc.prefix + name
	return &Queue{
		service:         svc,
		name:            name,
		nameWithPrefix:  queueName,
		url:             pointers.String(url),
		autoDel:         false,
		expire:          defaultExpireSecond,
		waitTimeSeconds: defaultWaitTimeSeconds,
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

// SetWaitTimeSeconds sets wait time timeout for message.
// Setting this value allows for a long polling workflow.
func (q *Queue) SetWaitTimeSeconds(sec int) {
	q.waitTimeSeconds = sec
}

// AddMessage adds message to the send spool.
// This assumes a Standard SQS Queue and not a FifoQueue
func (q *Queue) AddMessage(message string) {
	q.sendSpoolMu.Lock()
	defer q.sendSpoolMu.Unlock()

	num := fmt.Sprint(len(q.sendSpool) + 1)
	m := &SDK.SendMessageBatchRequestEntry{
		MessageBody: pointers.String(message),
		Id:          pointers.String(defaultMessageIDPrefix + num), // serial numbering for convenience sake
	}
	q.sendSpool = append(q.sendSpool, m)
}

// AddMessageWithGroupID adds a message to the send spool but adds the required attributes
// for a SQS FIFO Queue. This assumes the SQS FIFO Queue has ContentBasedDeduplication enabled.
func (q *Queue) AddMessageWithGroupID(message string, messageGroupID string) {
	q.sendSpoolMu.Lock()
	defer q.sendSpoolMu.Unlock()

	num := fmt.Sprint(len(q.sendSpool) + 1)
	m := &SDK.SendMessageBatchRequestEntry{
		MessageBody:    pointers.String(message),
		Id:             pointers.String(defaultMessageIDPrefix + num), // serial numbering for convenience sake
		MessageGroupId: pointers.String(messageGroupID),
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

	errList := newErrors()
	// send message
	for i := 0; i < len(messages); i++ {
		err := q.send(messages[i])
		if err != nil {
			q.service.Errorf("error on `SendMessageBatch` operation; queue=%s; error=%s;", q.nameWithPrefix, err.Error())
			errList.Add(err)
		}
	}
	q.sendSpool = nil

	if errList.HasError() {
		return errList
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

	wait := q.waitTimeSeconds

	if wait == 0 && num > 1 {
		wait = 1 // use long-polling for 1sec when to get multiple messages
	}

	// receive message from AWS api
	resp, err := q.service.client.ReceiveMessage(&SDK.ReceiveMessageInput{
		QueueUrl:            q.url,
		WaitTimeSeconds:     pointers.Long(wait),
		MaxNumberOfMessages: pointers.Long(num),
		VisibilityTimeout:   pointers.Long(q.expire),
	})
	if err != nil {
		q.service.Errorf("error on `ReceiveMessage` operation; queue=%s; error=%s;", q.nameWithPrefix, err.Error())
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

// ChangeMessageVisibility sends the request to AWS api to change visibility of the message.
func (q *Queue) ChangeMessageVisibility(msg *Message, timeoutInSeconds int) error {
	_, err := q.service.client.ChangeMessageVisibility(&SDK.ChangeMessageVisibilityInput{
		QueueUrl:          q.url,
		VisibilityTimeout: pointers.Long(timeoutInSeconds),
		ReceiptHandle:     msg.GetReceiptHandle(),
	})
	if err != nil {
		q.service.Errorf("error on `ChangeMessageVisibility`; queue=%s; error=%s;", q.nameWithPrefix, err.Error())
	}
	return err
}

// DeleteMessage sends the request to AWS api to delete the message.
func (q *Queue) DeleteMessage(msg *Message) error {
	_, err := q.service.client.DeleteMessage(&SDK.DeleteMessageInput{
		QueueUrl:      q.url,
		ReceiptHandle: msg.GetReceiptHandle(),
	})
	if err != nil {
		q.service.Errorf("error on `DeleteMessage`; queue=%s; error=%s;", q.nameWithPrefix, err.Error())
	}
	return err
}

// DeleteMessageWithReceipt sends the request to AWS api to delete the message.
func (q *Queue) DeleteMessageWithReceipt(msgReceipt string) error {
	_, err := q.service.client.DeleteMessage(&SDK.DeleteMessageInput{
		QueueUrl:      q.url,
		ReceiptHandle: pointers.String(msgReceipt),
	})
	if err != nil {
		q.service.Errorf("error on `DeleteMessage`; queue=%s; error=%s;", q.nameWithPrefix, err.Error())
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
	errList := newErrors()
	for i := 0; i < len(messages); i++ {
		err := q.delete(messages[i])
		if err != nil {
			errList.Add(err)
		}
	}
	q.deleteSpool = nil

	if errList.HasError() {
		return errList
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
		q.service.Errorf("error on `DeleteMessageBatch`; queue=%s; error=%s;", q.nameWithPrefix, err.Error())
		q.failedDelete = append(q.failedDelete, res.Failed...)
	}

	return err
}

// CountMessage sends request to AWS api to counts left messages in the Queue.
func (q *Queue) CountMessage() (visible int, invisible int, err error) {
	attr, err := q.service.GetQueueAttributes(*q.url,
		AttributeApproximateNumberOfMessages,
		AttributeApproximateNumberOfMessagesNotVisible,
	)
	if err != nil {
		return 0, 0, err
	}

	return attr.ApproximateNumberOfMessages, attr.ApproximateNumberOfMessagesNotVisible, nil
}

// GetAttributes sends request to AWS api to get the queue's attributes.
// `AttributeNames` will be set as `All`.
func (q *Queue) GetAttributes() (AttributesResponse, error) {
	return q.service.GetQueueAttributes(*q.url)
}

// Purge deletes all messages in the Queue.
func (q *Queue) Purge() error {
	_, err := q.service.client.PurgeQueue(&SDK.PurgeQueueInput{
		QueueUrl: q.url,
	})
	if err != nil {
		q.service.Errorf("error on `PurgeQueue` operation; queue=%s; error=%s;", q.nameWithPrefix, err.Error())
		return err
	}

	q.service.Infof("success on `PurgeQueue` operation; queue=%s;", q.nameWithPrefix)
	return nil
}
