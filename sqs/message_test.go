package sqs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	setTestEnv()
	createQueue("test")
}

func TestMessageString(t *testing.T) {
	assert := assert.New(t)
	svc := NewClient()
	q, _ := svc.GetQueue("test")

	cleanQueue(q)

	// prepare
	addTestMessage(q, 3)
	msg, err := q.FetchOne()
	assert.Nil(err)

	// test this feature
	str := msg.String()
	assert.Contains(str, `Body: "`+msg.Body()+`"`)
	assert.Contains(str, `MD5OfBody: "`)
	assert.Contains(str, `MessageId: "`+*msg.GetMessageID()+`"`)
	assert.Contains(str, `ReceiptHandle: "`+*msg.GetReceiptHandle()+`"`)
	cleanQueue(q)
}
