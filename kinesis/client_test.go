package kinesis

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

const defaultEndpoint = "http://localhost:4577"
const testStreamName = "test-stream"

func getTestConfig() config.Config {
	return config.Config{
		AccessKey: "access",
		SecretKey: "secret",
		Endpoint:  defaultEndpoint,
	}
}

func getTestClient(t *testing.T) *Kinesis {
	svc, err := New(getTestConfig())
	if err != nil {
		t.Errorf("error on create client; error=%s;", err.Error())
		t.FailNow()
	}
	return svc
}

func getTestStream(t *testing.T) *Stream {
	svc := getTestClient(t)
	svc.CreateStreamWithName(testStreamName)
	s, err := svc.GetStream(testStreamName)
	if err != nil {
		t.Errorf("error on getting test stream; error=%s;", err.Error())
		t.FailNow()
	}
	return s
}

func recreateTestStream(t *testing.T) {
	svc := getTestClient(t)
	svc.ForceDeleteStream(testStreamName)
	time.Sleep(20 * time.Millisecond)
	svc.CreateStreamWithName(testStreamName)
	time.Sleep(20 * time.Millisecond)
}

func TestNew(t *testing.T) {
	assert := assert.New(t)

	svc, err := New(getTestConfig())
	assert.NoError(err)
	assert.NotNil(svc.client)
	assert.Equal("kinesis", svc.client.ServiceName)
	assert.Equal(defaultEndpoint, svc.client.Endpoint)

	region := "us-west-1"
	svc, err = New(config.Config{
		Region: region,
	})
	assert.NoError(err)
	expectedEndpoint := "https://kinesis." + region + ".amazonaws.com"
	assert.Equal(expectedEndpoint, svc.client.Endpoint)
}

func TestSetLogger(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)
	assert.Equal(log.DefaultLogger, svc.logger)

	stdLogger := &log.StdLogger{}
	svc.SetLogger(stdLogger)
	assert.Equal(stdLogger, svc.logger)
}

func TestGetStream(t *testing.T) {
	assert := assert.New(t)
	recreateTestStream(t)
	svc := getTestClient(t)

	s, err := svc.GetStream(testStreamName)
	assert.NoError(err)
	assert.NotNil(s)

	svc.ForceDeleteStream("non_exist")
	s, err = svc.GetStream("non_exist")
	assert.Error(err)
	assert.Nil(s)

	// cache
	svc.streams["non_exist"] = svc.streams[testStreamName]
	s, err = svc.GetStream("non_exist")
	assert.NoError(err)
	assert.NotNil(s)
}

func TestCreateStreamWithName(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	// not exitst
	_ = svc.ForceDeleteStream("test2")
	time.Sleep(20 * time.Millisecond)
	has, err := svc.IsExistStream("test2")
	assert.NoError(err)
	assert.False(has)

	// create
	err = svc.CreateStreamWithName("test2")
	assert.NoError(err)

	// created
	has, err = svc.IsExistStream("test2")
	assert.NoError(err)
	assert.True(has)
}

func TestIsExistStream(t *testing.T) {
	assert := assert.New(t)
	_ = getTestStream(t)
	svc := getTestClient(t)

	has, err := svc.IsExistStream(testStreamName)
	assert.NoError(err)
	assert.True(has)

	// not exitst
	has, err = svc.IsExistStream("non-exitst-stream")
	assert.NoError(err)
	assert.False(has)
}
