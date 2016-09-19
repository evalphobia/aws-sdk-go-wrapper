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
	time.Sleep(50 * time.Millisecond)
	q.Purge()
	time.Sleep(50 * time.Millisecond)
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
	assert := assert.New(t)
	svc := getTestClient(t)
	q, _ := svc.GetQueue("test")

	assert.Equal(false, q.autoDel)
	q.AutoDelete(true)
	assert.Equal(true, q.autoDel)
	q.AutoDelete(false)
	assert.Equal(false, q.autoDel)
}

func TestSetExpire(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)
	q, _ := svc.GetQueue("test")

	assert.Equal(defaultExpireSecond, q.expire)
	q.SetExpire(10)
	assert.Equal(10, q.expire)
}

func TestAddMessage(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)
	q, _ := svc.GetQueue("test")

	q.AddMessage("foo msg")
	assert.Equal(1, len(q.sendSpool))
	msg := *(q.sendSpool[0].MessageBody)
	assert.Equal("foo msg", msg)
}

func TestAddMessageMap(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)
	q, _ := svc.GetQueue("test")

	m := make(map[string]interface{})
	m["number"] = 99
	m["title"] = "foo title"
	jsonMsg := `{"number":99,"title":"foo title"}`

	q.AddMessageMap(m)
	assert.Equal(1, len(q.sendSpool))
	msg := *(q.sendSpool[0].MessageBody)
	assert.Equal(jsonMsg, msg)
}

func TestSend(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)
	q, _ := svc.GetQueue("test")

	q.AddMessage("foo send")
	err := q.Send()
	assert.Nil(err)
}

func TestFetch(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)
	q, _ := svc.GetQueue("test")

	// prepare
	cleanQueue(q)
	num, _, _ := q.CountMessage()
	assert.Equal(0, num)

	// test this feature
	q.AutoDelete(true)
	addTestMessage(q, 3)
	list, err := q.Fetch(10)
	assert.Nil(err)
	assert.Equal(true, len(list) > 0)

	cleanQueue(q)
}

func TestFetchOne(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)
	q, _ := svc.GetQueue("test")

	// prepare
	cleanQueue(q)
	num, _, _ := q.CountMessage()
	assert.Equal(0, num)

	// test empty
	res, err := q.FetchOne()
	assert.Nil(err)
	assert.Nil(res)

	// test this feature
	addTestMessage(q, 3)
	res, err = q.FetchOne()
	assert.Nil(err)
	assert.Contains(res.Body(), "addTestMessage")

	cleanQueue(q)
}

func TestFetchBody(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)
	q, _ := svc.GetQueue("test")

	// prepare
	cleanQueue(q)
	num, _, _ := q.CountMessage()
	assert.Equal(0, num)

	// test empty
	res := q.FetchBody(10)
	assert.Equal(0, len(res))

	// test this feature
	q.AutoDelete(true)
	addTestMessage(q, 3)
	res = q.FetchBody(10)
	assert.True(len(res) > 0)
	assert.Contains(res[0], "addTestMessage")

	cleanQueue(q)
}

func TestFetchBodyOne(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)
	q, _ := svc.GetQueue("test")

	// prepare
	cleanQueue(q)
	num, _, _ := q.CountMessage()
	assert.Equal(0, num)

	// test empty
	res := q.FetchBodyOne()
	assert.Contains(res, "")

	// test this feature
	addTestMessage(q, 3)
	res = q.FetchBodyOne()
	assert.Contains(res, "addTestMessage")

	cleanQueue(q)
}

func TestAddDeleteList(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)
	q, _ := svc.GetQueue("test")
	assert.Equal(0, len(q.deleteSpool))

	sdkmsg := &SDK.Message{
		MessageId:     String("id"),
		ReceiptHandle: String("handle"),
	}

	// add single SDK.Message
	q.AddDeleteList(sdkmsg)
	assert.Equal(1, len(q.deleteSpool))

	// add slice SDK.Message
	q.AddDeleteList([]*SDK.Message{sdkmsg, sdkmsg})
	assert.Equal(3, len(q.deleteSpool))

	msg := &Message{sdkmsg}

	// add single message
	q.AddDeleteList(msg)
	assert.Equal(4, len(q.deleteSpool))

	// add slice message
	q.AddDeleteList([]*Message{msg, msg, msg})
	assert.Equal(7, len(q.deleteSpool))

}

func TestDeleteMessage(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)
	q, _ := svc.GetQueue("test")
	cleanQueue(q)

	// prepare messages
	addTestMessage(q, 3)
	msg, err := q.FetchOne()
	assert.Nil(err)

	// test this feature
	err = q.DeleteMessage(msg)
	assert.Nil(err)

	cleanQueue(q)
}

func TestDeleteListItems(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)
	q, _ := svc.GetQueue("test")
	cleanQueue(q)

	// prepare messages
	addTestMessage(q, 3)
	res, _ := q.Fetch(10)

	// test this feature
	q.AddDeleteList(res)
	err := q.DeleteListItems()
	assert.Nil(err)

	cleanQueue(q)

	// no message
	err = q.DeleteListItems()
	assert.Nil(err)

	cleanQueue(q)

	// over 10+ message
	addTestMessage(q, 30)
	var list []*Message
	for len(list) < 20 {
		res, _ := q.Fetch(10)
		list = append(list, res...)
	}
	q.AddDeleteList(list)
	err = q.DeleteListItems()
	assert.Nil(err)

	cleanQueue(q)
}

func TestCountMessage(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)
	q, _ := svc.GetQueue("test")
	cleanQueue(q)
	addTestMessage(q, 3)

	visible, invisible, err := q.CountMessage()
	sum := visible + invisible
	assert.Nil(err)
	assert.True(sum > 0)

	// test for increase
	addTestMessage(q, 3)
	visible2, invisible2, err := q.CountMessage()
	sum2 := visible2 + invisible2
	assert.Nil(err)
	assert.True(sum2 > sum)

	cleanQueue(q)
}

func TestPurge(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	// prepare message
	q, _ := svc.GetQueue("test")
	cleanQueue(q)
	addTestMessage(q, 3)

	visible, invisible, _ := q.CountMessage()
	sum := visible + invisible

	// test this feature
	err := q.Purge()
	assert.Nil(err)

	// make sure deleted
	assert.NotEqual(0, sum)
	visible2, invisible2, _ := q.CountMessage()
	assert.Equal(0, visible2)
	assert.Equal(0, invisible2)
}
