// SNS client

package sns

import (
	"encoding/json"
	"errors"

	AWS "github.com/awslabs/aws-sdk-go/aws"
	SNS "github.com/awslabs/aws-sdk-go/gen/sns"

	"github.com/evalphobia/aws-sdk-go-wrapper/auth"
	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

const (
	snsConfigSectionName = "sns"
	defaultRegion        = "us-west-1"
	defaultPrefix        = "dev_"

	AppTypeAPNS        = "apns"
	AppTypeAPNSSandbox = "apns_sandbox"
	AppTypeGCM         = "gcm"
)

var isProduction bool

type AmazonSNS struct {
	apps   map[string]*SNSApp
	topics map[string]*SNSTopic
	Client *SNS.SNS
}

// Create new AmazonSQS struct
func NewClient() *AmazonSNS {
	s := &AmazonSNS{}
	s.apps = make(map[string]*SNSApp)
	s.topics = make(map[string]*SNSTopic)
	region := config.GetConfigValue(snsConfigSectionName, "region", defaultRegion)
	s.Client = SNS.New(auth.Auth(), region, nil)
	if config.GetConfigValue(snsConfigSectionName, "app.production", "false") != "false" {
		isProduction = true
	} else {
		isProduction = false
	}
	return s
}

// Get SNSApp struct
func (s *AmazonSNS) GetApp(typ string) (*SNSApp, error) {
	// get the app from cache
	app, ok := s.apps[typ]
	if ok {
		return app, nil
	}
	arn := config.GetConfigValue(snsConfigSectionName, "app."+typ, "")
	if arn == "" {
		errMsg := "[SNS] error, cannot find ARN setting"
		log.Error(errMsg, typ)
		return nil, errors.New(errMsg)
	}
	app = &SNSApp{
		arn:      arn,
		platform: typ,
		client:   s,
	}
	s.apps[typ] = app
	return app, nil
}

// Get SNSApp struct of Apple Push Notification Service for iOS
func (s *AmazonSNS) GetAppAPNS() (*SNSApp, error) {
	if isProduction {
		return s.GetApp(AppTypeAPNS)
	} else {
		return s.GetApp(AppTypeAPNSSandbox)
	}
}

// Get SNSApp struct for Google Cloud Messaging for Android
func (s *AmazonSNS) GetAppGCM() (*SNSApp, error) {
	return s.GetApp(AppTypeGCM)
}

// Create SNS Topic and return `TopicARN`
func (s *AmazonSNS) createTopic(name string) (string, error) {
	prefix := config.GetConfigValue(snsConfigSectionName, "prefix", defaultPrefix)
	in := &SNS.CreateTopicInput{
		Name: AWS.String(prefix + name),
	}
	resp, err := s.Client.CreateTopic(in)
	if err != nil {
		log.Error("[SNS] error on `CreateTopic` operation, name="+name, err.Error())
		return "", err
	}
	return *resp.TopicARN, nil
}

// Create SNS Topic and return `TopicARN`
func (s *AmazonSNS) CreateTopic(name string) (*SNSTopic, error) {
	arn, err := s.createTopic(name)
	if err != nil {
		return nil, err
	}
	topic := &SNSTopic{
		name:   name,
		arn:    arn,
		client: s,
		sound:  "default",
	}
	return topic, nil
}

// Publish notification for arn(boty topic or endpoint)
func (s *AmazonSNS) Publish(arn string, msg string, opt map[string]interface{}) error {
	m := make(map[string]string)
	m["default"] = msg
	m["GCM"] = composeMessageGCM(msg)
	m["APNS"] = composeMessageAPNS(msg, opt)
	m["APNS_SANDBOX"] = m["APNS"]
	jsonString, _ := json.Marshal(m)
	resp, err := s.Client.Publish(&SNS.PublishInput{
		TargetARN:        AWS.String(arn),
		Message:          AWS.String(string(jsonString)),
		MessageStructure: AWS.String("json"),
	})
	if err != nil {
		log.Error("[SNS] error on `Publish` operation, arn="+arn, err.Error())
		return err
	}
	log.Info("[SNS] publish message", *resp.MessageID)
	return nil
}
