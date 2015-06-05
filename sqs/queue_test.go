package sqs

import (
	"fmt"
	"testing"
	"time"

	SDK "github.com/aws/aws-sdk-go/service/sqs"

	"github.com/stretchr/testify/assert"
)

// delete all messages from the Queue
func cleanQueue(q *Queue) {
	// purge has limitation for 60sec interval
	q.AutoDelete(true)
	time.Sleep(100 * time.Millisecond)
	for {
		num, num2, err := q.CountMessage()
		num += num2
		if num == 0 || err != nil {
			q.AutoDelete(false)
			time.Sleep(100 * time.Millisecond)
			return
		}
		q.Fetch(10)
		time.Sleep(500 * time.Millisecond)
	}
}

// add messages to the Queue
func addTestMessage(q *Queue, num int) {
	for i := 0; i < num; i++ {
		q.AddMessage(fmt.Sprintf("addTestMessage %d", i))
	}
	q.Send()
	time.Sleep(200 * time.Millisecond)
}

func TestAutoDelete(t *testing.T) {
	setTestEnv()
	createQueue("test")

	svc := NewClient()
	q, _ := svc.GetQueue("test")

	assert.Equal(t, false, q.autoDel)
	q.AutoDelete(true)
	assert.Equal(t, true, q.autoDel)
	q.AutoDelete(false)
	assert.Equal(t, false, q.autoDel)
}

func TestSetExpire(t *testing.T) {
	setTestEnv()
	createQueue("test")

	svc := NewClient()
	q, _ := svc.GetQueue("test")

	assert.Equal(t, defaultExpireSecond, q.expire)
	q.SetExpire(10)
	assert.Equal(t, 10, q.expire)
}

func TestAddMessage(t *testing.T) {
	setTestEnv()
	createQueue("test")

	svc := NewClient()
	q, _ := svc.GetQueue("test")

	q.AddMessage("foo msg")
	assert.Equal(t, 1, len(q.messages))
	msg := *(q.messages[0].MessageBody)
	assert.Equal(t, "foo msg", msg)
}

func TestAddMessageMap(t *testing.T) {
	setTestEnv()
	createQueue("test")

	svc := NewClient()
	q, _ := svc.GetQueue("test")

	m := make(map[string]interface{})
	m["number"] = 99
	m["title"] = "foo title"
	jsonMsg := `{"number":99,"title":"foo title"}`

	q.AddMessageMap(m)
	assert.Equal(t, 1, len(q.messages))
	msg := *(q.messages[0].MessageBody)
	assert.Equal(t, jsonMsg, msg)
}

func TestSend(t *testing.T) {
	setTestEnv()
	createQueue("test")

	svc := NewClient()
	q, _ := svc.GetQueue("test")

	q.AddMessage("foo send")
	err := q.Send()
	assert.Nil(t, err)
}

func TestFetch(t *testing.T) {
	setTestEnv()
	createQueue("test")

	svc := NewClient()
	q, _ := svc.GetQueue("test")

	// prepare
	cleanQueue(q)
	num, _, _ := q.CountMessage()
	assert.Equal(t, 0, num)

	// test this feature
	q.AutoDelete(true)
	addTestMessage(q, 3)
	res, err := q.Fetch(10)
	assert.Nil(t, err)
	assert.Equal(t, true, len(res.Messages) > 0)

	cleanQueue(q)
}

func TestFetchOne(t *testing.T) {
	setTestEnv()
	createQueue("test")

	svc := NewClient()
	q, _ := svc.GetQueue("test")

	// prepare
	cleanQueue(q)
	num, _, _ := q.CountMessage()
	assert.Equal(t, 0, num)

	// test empty
	res, err := q.FetchOne()
	assert.Nil(t, err)
	assert.Nil(t, res)

	// test this feature
	addTestMessage(q, 3)
	res, err = q.FetchOne()
	assert.Nil(t, err)
	assert.Contains(t, *res.Body, "addTestMessage")

	cleanQueue(q)
}

func TestFetchBody(t *testing.T) {
	setTestEnv()
	createQueue("test")

	svc := NewClient()
	q, _ := svc.GetQueue("test")

	// prepare
	cleanQueue(q)
	num, _, _ := q.CountMessage()
	assert.Equal(t, 0, num)

	// test empty
	res := q.FetchBody(10)
	assert.Equal(t, 0, len(res))

	// test this feature
	q.AutoDelete(true)
	addTestMessage(q, 3)
	res = q.FetchBody(10)
	assert.True(t, len(res) > 0)
	assert.Contains(t, res[0], "addTestMessage")

	cleanQueue(q)
}

func TestFetchBodyOne(t *testing.T) {
	setTestEnv()
	createQueue("test")

	svc := NewClient()
	q, _ := svc.GetQueue("test")

	// prepare
	cleanQueue(q)
	num, _, _ := q.CountMessage()
	assert.Equal(t, 0, num)

	// test empty
	res := q.FetchBodyOne()
	assert.Contains(t, res, "")

	// test this feature
	addTestMessage(q, 3)
	res = q.FetchBodyOne()
	assert.Contains(t, res, "addTestMessage")

	cleanQueue(q)
}

func TestAddDeleteList(t *testing.T) {
	setTestEnv()
	createQueue("test")

	svc := NewClient()
	q, _ := svc.GetQueue("test")
	assert.Equal(t, 0, len(q.delMessages))

	msg := &SDK.Message{
		MessageID:     String(""),
		ReceiptHandle: String(""),
	}

	// add single message
	q.AddDeleteList(msg)
	assert.Equal(t, 1, len(q.delMessages))

	// add slice message
	q.AddDeleteList([]*SDK.Message{msg, msg})
	assert.Equal(t, 3, len(q.delMessages))
}

func TestDeleteMessage(t *testing.T) {
	setTestEnv()
	createQueue("test")

	svc := NewClient()
	q, _ := svc.GetQueue("test")
	cleanQueue(q)

	// prepare messages
	addTestMessage(q, 3)
	msg, err := q.FetchOne()
	assert.Nil(t, err)

	// test this feature
	err = q.DeleteMessage(msg)
	assert.Nil(t, err)

	cleanQueue(q)
}

func TestDeleteListItems(t *testing.T) {
	setTestEnv()
	createQueue("test")

	svc := NewClient()
	q, _ := svc.GetQueue("test")
	cleanQueue(q)

	// prepare messages
	addTestMessage(q, 3)
	res, _ := q.Fetch(10)

	// test this feature
	q.AddDeleteList(res)
	err := q.DeleteListItems()
	assert.Nil(t, err)

	cleanQueue(q)
}

func TestCountMessage(t *testing.T) {
	setTestEnv()
	createQueue("test")

	svc := NewClient()
	q, _ := svc.GetQueue("test")
	cleanQueue(q)
	addTestMessage(q, 3)

	visible, invisible, err := q.CountMessage()
	sum := visible + invisible
	assert.Nil(t, err)
	assert.True(t, sum > 0)

	// test for increase
	addTestMessage(q, 3)
	visible2, invisible2, err := q.CountMessage()
	sum2 := visible2 + invisible2
	assert.Nil(t, err)
	assert.True(t, sum2 > sum)

	cleanQueue(q)
}

func TestPurge(t *testing.T) {
	setTestEnv()
	createQueue("test")

	svc := NewClient()
	if svc.client.Endpoint == defaultEndpoint {
		t.Skip("fakesqs does not implement Purge() yet.")
	}

	// prepare message
	q, _ := svc.GetQueue("test")
	cleanQueue(q)
	addTestMessage(q, 3)

	visible, invisible, _ := q.CountMessage()
	sum := visible + invisible

	// test this feature
	err := q.Purge()
	assert.Nil(t, err)

	// make sure deleted
	assert.NotEqual(t, 0, sum)
	visible2, invisible2, _ := q.CountMessage()
	assert.Equal(t, 0, visible2)
	assert.Equal(t, 0, invisible2)
}
