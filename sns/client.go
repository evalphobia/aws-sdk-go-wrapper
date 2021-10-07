package sns

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/aws/aws-sdk-go/aws/session"
	SDK "github.com/aws/aws-sdk-go/service/sns"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

// Application types
const (
	AppTypeAPNS        = "APNS"
	AppTypeAPNSSandbox = "APNS_SANDBOX"
	AppTypeGCM         = "GCM"

	ProtocolApplication = "application"
)

const (
	topicMaxDeviceNumber = 10000
	messageBodyLimit     = 2000
)

// SNS is AWS SNS client and has platform application and topic list.
type SNS struct {
	client *SDK.SNS

	logger       log.Logger
	prefix       string
	isProduction bool
	appsMu       sync.RWMutex
	apps         map[string]*PlatformApplication

	platforms Platforms
}

// New returns initialized *SNS.
// use ARNs in given Platforms.
func New(conf config.Config, pf Platforms) (*SNS, error) {
	sess, err := conf.Session()
	if err != nil {
		return nil, err
	}

	svc := NewFromSession(sess)
	svc.prefix = conf.DefaultPrefix
	svc.platforms = pf
	return svc, nil
}

// NewFromSession returns initialized *SNS from aws.Session.
func NewFromSession(sess *session.Session) *SNS {
	return &SNS{
		client: SDK.New(sess),
		logger: log.DefaultLogger,
		apps:   make(map[string]*PlatformApplication),
	}
}

// GetClient gets aws client.
func (svc *SNS) GetClient() *SDK.SNS {
	return svc.client
}

// SetLogger sets logger.
func (svc *SNS) SetLogger(logger log.Logger) {
	svc.logger = logger
}

// SetPrefix sets prefix.
func (svc *SNS) SetPrefix(prefix string) {
	svc.prefix = prefix
}

// SetPlatforms sets platforms.
func (svc *SNS) SetPlatforms(pf Platforms) {
	svc.platforms = pf
}

// SetAsProduction sets productiom mode flag.
func (svc *SNS) SetAsProduction() {
	svc.platforms.Production = true
}

// ================================
// PlatformApplication Operations
// ================================

// GetPlatformApplicationApple returns *PlatformApplication for iOS APNS.
func (svc *SNS) GetPlatformApplicationApple() *PlatformApplication {
	return svc.getPlatformApplicationByType(AppTypeAPNS)
}

// GetPlatformApplicationGoogle returns *PlatformApplication for Android GCM.
func (svc *SNS) GetPlatformApplicationGoogle() *PlatformApplication {
	return svc.getPlatformApplicationByType(AppTypeGCM)
}

// getPlatformApplicationByType returns *PlatformApplication.
func (svc *SNS) getPlatformApplicationByType(typ string) *PlatformApplication {
	typ = strings.ToUpper(typ)

	// use apns sandbox when it's not in production env
	if typ == AppTypeAPNS && !svc.isProduction {
		typ = AppTypeAPNSSandbox
	}

	// get the app from cache
	svc.appsMu.RLock()
	app, ok := svc.apps[typ]
	svc.appsMu.RUnlock()
	if ok {
		return app
	}

	app = svc.newPlatformApplication(svc.platforms.GetARNByType(typ), typ)
	svc.appsMu.Lock()
	svc.apps[typ] = app
	svc.appsMu.Unlock()
	return app
}

func (svc *SNS) newPlatformApplication(arn, pf string) *PlatformApplication {
	return &PlatformApplication{
		svc:      svc,
		arn:      arn,
		platform: pf,
	}
}

// ================================
// Topic Operations
// ================================

// CreateTopic creates Topic.
func (svc *SNS) CreateTopic(name string) (*Topic, error) {
	arn, err := svc.createTopic(name)
	if err != nil {
		return nil, err
	}

	topic := NewTopic(svc, arn, name)
	return topic, nil
}

// createTopic operates CreateTopic and return `TopicARN`.
func (svc *SNS) createTopic(name string) (topicARN string, err error) {
	topicName := svc.prefix + name
	in := &SDK.CreateTopicInput{
		Name: pointers.String(topicName),
	}
	resp, err := svc.client.CreateTopic(in)
	if err != nil {
		svc.Errorf("error on `CreateTopic` operation; name=%s; error=%s;", name, err.Error())
		return "", err
	}
	return *resp.TopicArn, nil
}

// Publish sends mobile notifications to the ARN (topic or endpoint).
func (svc *SNS) Publish(arn string, msg string, options map[string]interface{}) error {
	// trim message size
	msg = truncateMessage(msg)

	m := make(map[string]string)
	m["default"] = msg

	// GCM
	var err error
	m[AppTypeGCM], err = composeMessageGCM(msg, options)
	if err != nil {
		svc.Errorf("error on composeMessageGCM; msg=%s; error=%s;", msg, err.Error())
		return err
	}

	// APNS
	switch {
	case svc.platforms.Production:
		m[AppTypeAPNS], err = composeMessageAPNS(msg, options)
	default:
		m[AppTypeAPNSSandbox], err = composeMessageAPNS(msg, options)
	}
	if err != nil {
		svc.Errorf("error on composeMessageAPNS; msg=%s; error=%s;", msg, err.Error())
		return err
	}

	jsonByte, err := json.Marshal(m)
	if err != nil {
		svc.Errorf("error on json.Marshal; arn=%s; error=%s;", arn, err.Error())
		return err
	}

	_, err = svc.client.Publish(&SDK.PublishInput{
		TargetArn:        pointers.String(arn),
		Message:          pointers.String(string(jsonByte)),
		MessageStructure: pointers.String("json"),
	})
	if err != nil {
		svc.Errorf("error on `Publish` operation; arn=%s; error=%s;", arn, err.Error())
	}
	return err
}

// truncateMessage limits message size to the allowed payload size.
func truncateMessage(msg string) string {
	if len(msg) <= messageBodyLimit {
		return msg
	}

	runes := []rune(msg[:messageBodyLimit])
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

// PublishAPNSByToken sends push message for iOS device by device token
func (svc *SNS) PublishAPNSByToken(token string, msg string, badge int) error {
	return svc.PublishByToken(AppTypeAPNS, token, msg, badge)
}

// PublishGCMByToken sends push message for Android device by device token
func (svc *SNS) PublishGCMByToken(token string, msg string, badge int) error {
	return svc.PublishByToken(AppTypeGCM, token, msg, badge)
}

// PublishByToken sends push message by device token
func (svc *SNS) PublishByToken(device, token string, msg string, badge int) error {
	ep, err := svc.RegisterEndpoint(device, token)
	if err != nil {
		return err
	}
	return ep.Publish(msg, badge)
}

// BulkPublishByDevice sends mobile notification to many endpoints.
// (supports single device only)
func (svc *SNS) BulkPublishByDevice(device string, tokens []string, msg string) error {
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
			topic.Subscribe(e.arn, e.protocol)
			wg.Done()
		}(token)
	}
	wg.Wait()
	topic.Publish(msg)
	topic.Delete()
	return nil
}

// BulkPublish sends mobile notification for many endpoints.
// tokens is map of string slices, each key stands for device, like "android"/"ios"
// ex) tokens := map[string][]string{ "android": []string{"token1", "token2"}, "ios": []string{"token3", "token4"}}
func (svc *SNS) BulkPublish(tokens map[string][]string, msg string) error {
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

// RegisterEndpoint creates endpoint(device) to platform application.
func (svc *SNS) RegisterEndpoint(device, token string) (*PlatformEndpoint, error) {
	app, err := svc.getApp(device)
	if err != nil {
		return nil, err
	}
	return app.CreateEndpoint(token)
}

// RegisterEndpointWithUserData creates endpoint(device) and CustomUserData to platform application.
func (svc *SNS) RegisterEndpointWithUserData(device, token, userData string) (*PlatformEndpoint, error) {
	app, err := svc.getApp(device)
	if err != nil {
		return nil, err
	}
	return app.CreateEndpointWithUserData(token, userData)
}

func (svc *SNS) getApp(device string) (*PlatformApplication, error) {
	device = strings.ToUpper(device)

	switch device {
	case AppTypeAPNS, "IOS", "APPLE":
		return svc.GetPlatformApplicationApple(), nil
	case AppTypeGCM, "ANDROID", "GOOGLE":
		return svc.GetPlatformApplicationGoogle(), nil
	}

	err := fmt.Errorf("unsupported device")
	svc.Errorf("error getApp; device=%s; error=%s;", device, err.Error())
	return nil, err
}

// GetEndpoint gets *PlatformEndpoint by ARN.
func (svc *SNS) GetEndpoint(arn string) (*PlatformEndpoint, error) {
	in := &SDK.GetEndpointAttributesInput{
		EndpointArn: pointers.String(arn),
	}
	resp, err := svc.client.GetEndpointAttributes(in)
	if err != nil {
		svc.Errorf("error on `GetEndpointAttributes` operation; arn=%s; error=%s;", arn, err.Error())
		return nil, err
	}

	attr := resp.Attributes
	ep := svc.newApplicationEndpoint(arn)
	ep.token = *attr["Token"]
	ep.enable, err = strconv.ParseBool(*attr["Enabled"])
	if err != nil {
		svc.Errorf("error ParseBool(endpoint.Enabled); arn=%s; Enabled=%s; error=%s;", arn, *attr["Enabled"], err.Error())
	}
	return ep, err
}

// GetPlatformApplicationAttributes executes `GetPlatformApplicationAttributes`.
func (svc *SNS) GetPlatformApplicationAttributes(arn string) (PlatformAttributes, error) {
	resp, err := svc.client.GetPlatformApplicationAttributes(&SDK.GetPlatformApplicationAttributesInput{
		PlatformApplicationArn: pointers.String(arn),
	})
	if err != nil {
		svc.Errorf("error on `GetPlatformApplicationAttributes` operation; arn=%s; error=%s;", arn, err.Error())
		return PlatformAttributes{}, err
	}

	return NewPlatformAttributesFromMap(resp.Attributes), nil
}

func (svc *SNS) newApplicationEndpoint(arn string) *PlatformEndpoint {
	return &PlatformEndpoint{
		svc:      svc,
		arn:      arn,
		protocol: ProtocolApplication,
	}
}

// Infof logging information.
func (svc *SNS) Infof(format string, v ...interface{}) {
	svc.logger.Infof("SNS", format, v...)
}

// Errorf logging error information.
func (svc *SNS) Errorf(format string, v ...interface{}) {
	svc.logger.Errorf("SNS", format, v...)
}
