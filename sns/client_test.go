package sns

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
)

type testConfig struct {
	config map[string]interface{}
}

func (c *testConfig) GetConfigValue(sec, key, df string) string {
	v, ok := c.config[key]
	if !ok {
		return ""
	}
	return config.ParseToString(v)
}

func setTestConfig() {
	conf := make(map[string]interface{})
	conf["app.gcm"] = "arn:aws:sns:us-east-1:0000000000:app/GCM/foo_gcm"
	conf["app.apns"] = "arn:aws:sns:us-east-1:0000000000:app/GCM/foo_apns"
	conf["app.apns_sandbox"] = "arn:aws:sns:us-east-1:0000000000:app/GCM/foo_apns"
	conf["app.production"] = false
	c := &testConfig{conf}
	config.SetConfig(c)
}

func setTestEnv() {
	os.Clearenv()
	os.Setenv("AWS_ACCESS_KEY_ID", "access")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
}

func TestNewClient(t *testing.T) {
	setTestEnv()

	svc := NewClient()
	assert.NotNil(t, svc.Client)
	assert.Equal(t, "sns", svc.Client.ServiceName)
	assert.Equal(t, defaultEndpoint, svc.Client.Endpoint)

	region := "us-west-1"
	os.Setenv("AWS_REGION", region)
	svc2 := NewClient()
	endpoint := "https://sns." + region + ".amazonaws.com"
	assert.Equal(t, endpoint, svc2.Client.Endpoint)
}

func TestGetApp(t *testing.T) {
	setTestEnv()
	setTestConfig()

	svc := NewClient()
	app, err := svc.GetApp("gcm")
	assert.Nil(t, err)
	assert.NotNil(t, app)

	app, err = svc.GetApp("apns")
	assert.Nil(t, err)
	assert.NotNil(t, app)

	app, err = svc.GetApp("apns_sandbox")
	assert.Nil(t, err)
	assert.NotNil(t, app)

	app, err = svc.GetApp("foo")
	assert.NotNil(t, err)
	assert.Nil(t, app)
}

func TestGetAppAPNS(t *testing.T) {
	setTestEnv()
	setTestConfig()

	svc := NewClient()
	app, err := svc.GetAppAPNS()
	assert.Nil(t, err)
	assert.NotNil(t, app)
}

func TestGetAppGCM(t *testing.T) {
	setTestEnv()
	setTestConfig()

	svc := NewClient()
	app, err := svc.GetAppGCM()
	assert.Nil(t, err)
	assert.NotNil(t, app)
}

func TestCreateTopic(t *testing.T) {
	setTestEnv()
	setTestConfig()

	svc := NewClient()
	app, err := svc.CreateTopic("fooTopic")
	assert.Nil(t, err)
	assert.NotNil(t, app)
}

func TestPublish(t *testing.T) {
	setTestEnv()
	setTestConfig()

	svc := NewClient()
	opt := make(map[string]interface{})
	topic, _ := svc.CreateTopic("fooTopic")
	err := svc.Publish(topic.arn, "message", opt)

	t.Skip("fakesns does not implement Publish() yet.")
	_ = err
}

func TestTruncateMessage(t *testing.T) {
	str := "foobar"
	msg := truncateMessage(str)
	assert.Equal(t, str, msg)

	var bigStr string
	for i := 0; i < 3000; i++ {
		bigStr = bigStr + "a"
	}
	msg = truncateMessage(bigStr)
	assert.NotEqual(t, bigStr, msg)
	assert.Equal(t, 2000, len(msg))
}

func TestRegisterEndpoint(t *testing.T) {
	setTestEnv()
	setTestConfig()
	svc := NewClient()
	
	ep, err := svc.RegisterEndpoint("foo", "token")
	assert.NotNil(t, err)
	assert.Nil(t, ep)

	ep, err = svc.RegisterEndpoint("apns", "token")
	ep, err = svc.RegisterEndpoint("gcm", "token")
	t.Skip("fakesns does not implement CreatePlatformEndpoint() yet.")
}

func TestBulkPublishByDevice(t *testing.T) {
	setTestEnv()
	setTestConfig()
	svc := NewClient()
	
	err := svc.BulkPublishByDevice("ios", []string{"fooEndpoint"}, "message")
	assert.Nil(t, err)
}

func TestBulkPublish(t *testing.T) {
	setTestEnv()
	setTestConfig()
	svc := NewClient()
	
	tokens := map[string][]string{ "android": []string{"token1", "token2"}, "ios": []string{"token3", "token4"}}
	err := svc.BulkPublish(tokens, "message")
	assert.Nil(t, err)
}
