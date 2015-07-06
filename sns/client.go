// SNS client

package sns

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
	"unicode/utf8"

	SDK "github.com/aws/aws-sdk-go/service/sns"

	"github.com/evalphobia/aws-sdk-go-wrapper/auth"
	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

const (
	snsConfigSectionName = "sns"
	defaultRegion        = "us-west-1"
	defaultEndpoint      = "http://localhost:9292"
	defaultPrefix        = "dev_"

	AppTypeAPNS        = "apns"
	AppTypeAPNSSandbox = "apns_sandbox"
	AppTypeGCM         = "gcm"

	topicMaxDeviceNumber = 10000
	MessageBodyLimit     = 2000
)

var isProduction bool

type AmazonSNS struct {
	apps   map[string]*SNSApp
	topics map[string]*SNSTopic
	Client *SDK.SNS
}

// Create new AmazonSQS struct
func NewClient() *AmazonSNS {
	svc := &AmazonSNS{}
	svc.apps = make(map[string]*SNSApp)
	svc.topics = make(map[string]*SNSTopic)
	region := config.GetConfigValue(snsConfigSectionName, "region", auth.EnvRegion())
	awsConf := auth.NewConfig(region)
	endpoint := config.GetConfigValue(snsConfigSectionName, "endpoint", "")
	switch {
	case endpoint != "":
		awsConf.Endpoint = endpoint
	case region == "":
		awsConf.Region = defaultRegion
		awsConf.Endpoint = defaultEndpoint
	}
	svc.Client = SDK.New(awsConf)
	if config.GetConfigValue(snsConfigSectionName, "app.production", "false") != "false" {
		isProduction = true
	} else {
		isProduction = false
	}
	return svc
}

// Get SNSApp struct
func (svc *AmazonSNS) GetApp(typ string) (*SNSApp, error) {
	// get the app from cache
	app, ok := svc.apps[typ]
	if ok {
		return app, nil
	}
	arn := config.GetConfigValue(snsConfigSectionName, "app."+typ, "")
	if arn == "" {
		errMsg := "[SNS] error, cannot find ARN setting"
		log.Error(errMsg, typ)
		return nil, errors.New(errMsg)
	}
	app = NewApp(arn, typ, svc)
	svc.apps[typ] = app
	return app, nil
}

// Get SNSApp struct of Apple Push Notification Service for iOS
func (svc *AmazonSNS) GetAppAPNS() (*SNSApp, error) {
	if isProduction {
		return svc.GetApp(AppTypeAPNS)
	} else {
		return svc.GetApp(AppTypeAPNSSandbox)
	}
}

// Get SNSApp struct for Google Cloud Messaging for Android
func (svc *AmazonSNS) GetAppGCM() (*SNSApp, error) {
	return svc.GetApp(AppTypeGCM)
}

// Create SNS Topic and return `TopicARN`
func (svc *AmazonSNS) createTopic(name string) (string, error) {
	prefix := config.GetConfigValue(snsConfigSectionName, "prefix", defaultPrefix)
	in := &SDK.CreateTopicInput{
		Name: String(prefix + name),
	}
	resp, err := svc.Client.CreateTopic(in)
	if err != nil {
		log.Error("[SNS] error on `CreateTopic` operation, name="+name, err.Error())
		return "", err
	}
	return *resp.TopicARN, nil
}

// Create SNS Topic and return `TopicARN`
func (svc *AmazonSNS) CreateTopic(name string) (*SNSTopic, error) {
	arn, err := svc.createTopic(name)
	if err != nil {
		return nil, err
	}
	topic := NewTopic(arn, name, svc)
	return topic, nil
}

// Publish notification for arn(topic or endpoint)
func (svc *AmazonSNS) Publish(arn string, msg string, opt map[string]interface{}) error {
	msg = truncateMessage(msg)
	m := make(map[string]string)
	m["default"] = msg
	m["GCM"] = composeMessageGCM(msg)
	m["APNS"] = composeMessageAPNS(msg, opt)
	m["APNS_SANDBOX"] = m["APNS"]
	jsonString, _ := json.Marshal(m)
	resp, err := svc.Client.Publish(&SDK.PublishInput{
		TargetARN:        String(arn),
		Message:          String(string(jsonString)),
		MessageStructure: String("json"),
	})
	if err != nil {
		log.Error("[SNS] error on `Publish` operation, arn="+arn, err.Error())
		return err
	}
	log.Info("[SNS] publish message", *resp.MessageID)
	return nil
}

// Limit message size to the allowed payload size
func truncateMessage(msg string) string {
	if len(msg) <= MessageBodyLimit {
		return msg
	}
	runes := []rune(msg[:MessageBodyLimit])
	valid := len(runes)
	// traverse runes from last string and detect invalid string
	for i := valid; ; {
		i--
		if runes[i] != utf8.RuneError {
			break
		}
		valid = i
	}
	return string(runes[0:valid])
}

// Register endpoint(device) to application
func (svc *AmazonSNS) RegisterEndpoint(device, token string) (*SNSEndpoint, error) {
	var app *SNSApp
	var err error
	switch device {
	case "ios", "apns":
		app, err = svc.GetAppAPNS()
	case "android", "gcm":
		app, err = svc.GetAppGCM()
	default:
		errMsg := "[SNS] Unsupported device, device=" + device
		log.Error(errMsg, token)
		return nil, errors.New(errMsg)
	}
	if err != nil {
		return nil, err
	}
	return app.CreateEndpoint(token)
}

// PublishAPNSByToken sends push message for iOS device by device token
func (svc *AmazonSNS) PublishAPNSByToken(token string, msg string, badge int) error {
	return svc.PublishByToken(AppTypeAPNS, token, msg, badge)
}

// PublishGCMByToken sends push message for Android device by device token
func (svc *AmazonSNS) PublishGCMByToken(token string, msg string, badge int) error {
	return svc.PublishByToken(AppTypeGCM, token, msg, badge)
}

// PublishByToken sends push message by device token
func (svc *AmazonSNS) PublishByToken(device, token string, msg string, badge int) error {
	ep, err := svc.RegisterEndpoint(device, token)
	if err != nil {
		return err
	}
	return ep.Publish(msg, badge)
}

// Publish notification for many endpoints
// (supports single device only)
func (svc *AmazonSNS) BulkPublishByDevice(device string, tokens []string, msg string) error {
	name := fmt.Sprintf("%d", time.Now().UnixNano()) + "_" + device
	topic, err := svc.CreateTopic(name)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, token := range tokens {
		wg.Add(1)
		go func(t string) {
			e, err := svc.RegisterEndpoint(device, t)
			if err != nil {
				wg.Done()
				return
			}
			topic.Subscribe(e)
			wg.Done()
		}(token)
	}
	wg.Wait()
	topic.Publish(msg)
	topic.Delete()
	return nil
}

// Publish notification for many endpoints
// tokens is map of string slices, each key stands for device, like "android"/"ios"
// ex) tokens := map[string][]string{ "android": []string{"token1", "token2"}, "ios": []string{"token3", "token4"}}
func (svc *AmazonSNS) BulkPublish(tokens map[string][]string, msg string) error {
	for device, t := range tokens {
		l := len(t)
		max := (l-1)/topicMaxDeviceNumber + 1
		for i := 0; i < max; i++ {
			from := i * topicMaxDeviceNumber
			to := (i + 1) * topicMaxDeviceNumber
			if l < to {
				to = l
			}
			svc.BulkPublishByDevice(device, t[from:to], msg)
		}
	}
	return nil
}
