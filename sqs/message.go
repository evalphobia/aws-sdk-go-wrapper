// SQS Queue

package sqs

import (
	SDK "github.com/aws/aws-sdk-go/service/sqs"
)

// Message is SQS Message wrapper struct.
type Message struct {
	message *SDK.Message
}

// NewMessage returns initialized *Message.
func NewMessage(msg *SDK.Message) *Message {
	return &Message{msg}
}

func (m *Message) String() string {
	return m.message.String()
}

// Body returns message body.
func (m *Message) Body() string {
	return *m.message.Body
}

// GetMessageID returns pointer of message id.
func (m *Message) GetMessageID() *string {
	return m.message.MessageId
}

// GetReceiptHandle returns pointer of ReceiptHandle.
func (m *Message) GetReceiptHandle() *string {
	return m.message.ReceiptHandle
}
