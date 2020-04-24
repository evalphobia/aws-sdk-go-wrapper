package sns

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
)

const (
	testAppleARN    = "arn:aws:sns:us-east-1:0000000000:app/GCM/foo_apns"
	testGoogleARN   = "arn:aws:sns:us-east-1:0000000000:app/GCM/foo_gcm"
	defaultEndpoint = "http://localhost:9292"
)

func getTestConfig() config.Config {
	return config.Config{
		AccessKey: "access",
		SecretKey: "secret",
		Endpoint:  defaultEndpoint,
	}
}

func getTestClient(t *testing.T) *SNS {
	svc, err := New(getTestConfig(), Platforms{
		Production: false,
		Apple:      testAppleARN,
		Google:     testGoogleARN,
	})
	if err != nil {
		t.Errorf("error on create client; error=%s;", err.Error())
		t.FailNow()
	}
	return svc
}

func TestNewClient(t *testing.T) {
	assert := assert.New(t)

	pf := Platforms{
		Production: false,
		Apple:      "arn:aws:sns:us-east-1:0000000000:app/GCM/foo_apns",
		Google:     "arn:aws:sns:us-east-1:0000000000:app/GCM/foo_gcm",
	}

	svc, err := New(getTestConfig(), pf)
	assert.NoError(err)
	assert.NotNil(svc.client)
	assert.Equal("sns", svc.client.ServiceName)
	assert.Equal(defaultEndpoint, svc.client.Endpoint)

	region := "us-west-1"
	svc, err = New(config.Config{
		Region: region,
	}, pf)
	assert.NoError(err)
	expectedEndpoint := "https://sns." + region + ".amazonaws.com"
	assert.Equal(expectedEndpoint, svc.client.Endpoint)
}

func TestGetApp(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		typ              string
		expectedPlatform string
		expectedARN      string
	}{
		{"apns", AppTypeAPNSSandbox, testAppleARN},
		{"apns_sandbox", AppTypeAPNSSandbox, testAppleARN},
		{"gcm", AppTypeGCM, testGoogleARN},
		{"foo", "FOO", ""},
	}

	svc := getTestClient(t)
	for _, tt := range tests {
		target := fmt.Sprintf("%+v", tt)

		app := svc.getPlatformApplicationByType(tt.typ)
		assert.NotNil(app, target)
		assert.Equal(tt.expectedPlatform, app.platform, target)
		assert.Equal(tt.expectedARN, app.arn, target)
	}
}

func TestGetPlatformApplicationApple(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	app := svc.GetPlatformApplicationApple()
	assert.NotNil(app)
	assert.Equal(AppTypeAPNSSandbox, app.platform)
	assert.Equal(testAppleARN, app.arn)
}

func TestGetPlatformApplicationGoogle(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	app := svc.GetPlatformApplicationGoogle()
	assert.NotNil(app)
	assert.Equal(AppTypeGCM, app.platform)
	assert.Equal(testGoogleARN, app.arn)
}

func TestCreateTopic(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	topic, err := svc.CreateTopic("fooTopic")
	assert.Nil(err)
	assert.NotNil(topic)
	assert.Equal("fooTopic", topic.name)
	assert.Equal("default", topic.sound)
}

func TestPublish(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	if svc.client.Endpoint == defaultEndpoint {
		t.Skip("fakesns does not implement Publish() yet.")
	}

	opt := make(map[string]interface{})
	topic, _ := svc.CreateTopic("fooTopic")
	err := svc.Publish(topic.arn, "message", opt)
	assert.NoError(err)
}

func TestTruncateMessage(t *testing.T) {
	assert := assert.New(t)

	str := "foobar"
	msg := truncateMessage(str)
	assert.Equal(str, msg)

	var bigStr string
	for i := 0; i < 3000; i++ {
		bigStr = bigStr + "a"
	}
	msg = truncateMessage(bigStr)
	assert.NotEqual(bigStr, msg)
	assert.Equal(2000, len(msg))
}

func TestRegisterEndpoint(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	ep, err := svc.RegisterEndpoint("foo", "token")
	assert.NotNil(err)
	assert.Nil(ep)

	ep, err = svc.RegisterEndpoint("apns", "token")
	ep, err = svc.RegisterEndpoint("gcm", "token")
	t.Skip("fakesns does not implement CreatePlatformEndpoint() yet.")
}

func TestBulkPublishByDevice(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	err := svc.BulkPublishByDevice("ios", []string{"fooEndpoint"}, "message")
	assert.Nil(err)
}

func TestBulkPublish(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	tokens := map[string][]string{"android": {"token1", "token2"}, "ios": {"token3", "token4"}}
	err := svc.BulkPublish(tokens, "message")
	assert.Nil(err)
}

func TestGetPlatformApplicationAttributes(t *testing.T) {
	a := assert.New(t)
	svc := getTestClient(t)

	list := []string{
		testAppleARN,
		testGoogleARN,
	}

	t.Skip("fakesns does not implement GetPlatformApplicationAttributes() yet.")

	for _, v := range list {
		resp, err := svc.GetPlatformApplicationAttributes(v)
		a.NoError(err, v)
		a.True(resp.HasEnabled, v)
		a.True(resp.Enabled, v)
	}
}
