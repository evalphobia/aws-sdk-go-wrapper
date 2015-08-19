// SQS Queue

package sqs

import (
	SDK "github.com/aws/aws-sdk-go/service/sqs"
)

// SQS Message wrapper struct
type Message struct {
	message *SDK.Message
}

func NewMessage(msg *SDK.Message) *Message {
	return &Message{msg}
}

func (m Message) String() string {
	return m.message.String()
}

func (m Message) Body() string {
	return *m.message.Body
}

func (m Message) GetMessageID() *string {
	return m.message.MessageId
}

func (m Message) GetReceiptHandle() *string {
	return m.message.ReceiptHandle
}
