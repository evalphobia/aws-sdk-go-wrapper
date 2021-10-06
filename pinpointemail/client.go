package pinpointemail

import (
	"github.com/aws/aws-sdk-go/aws/session"
	SDK "github.com/aws/aws-sdk-go/service/pinpointemail"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

const (
	serviceName = "PinpointEmail"
)

// PinpointEmail has PinpointEmail client.
type PinpointEmail struct {
	client *SDK.PinpointEmail

	logger log.Logger
}

// New returns initialized *PinpointEmail.
func New(conf config.Config) (*PinpointEmail, error) {
	sess, err := conf.Session()
	if err != nil {
		return nil, err
	}

	return NewFromSession(sess), nil
}

// NewFromSession returns initialized *PinpointEmail from aws.Session.
func NewFromSession(sess *session.Session) *PinpointEmail {
	return &PinpointEmail{
		client: SDK.New(sess),
		logger: log.DefaultLogger,
	}
}

// SetLogger sets logger.
func (svc *PinpointEmail) SetLogger(logger log.Logger) {
	svc.logger = logger
}

// SendEmailSimple sends email from simple arguments.
func (svc *PinpointEmail) SendEmailSimple(subject, body, from string, to ...string) (string, error) {
	return svc.SendEmail(EmailInput{
		From: from,
		To:   to,
		Content: Content{
			Subject: subject,
			Body:    body,
		},
	})
}

// SendEmail sends email.
func (svc *PinpointEmail) SendEmail(in EmailInput) (string, error) {
	return svc.DoSendEmail(in.ToInput())
}

// DoSendEmail executes SendEmail operation.
func (svc *PinpointEmail) DoSendEmail(in *SDK.SendEmailInput) (string, error) {
	out, err := svc.client.SendEmail(in)
	if err != nil {
		svc.Errorf("error on `SendEmail` operation; error=%s;", err.Error())
		return "", err
	}

	id := *out.MessageId
	return id, nil
}

// Infof logging information.
func (svc *PinpointEmail) Infof(format string, v ...interface{}) {
	svc.logger.Infof(serviceName, format, v...)
}

// Errorf logging error information.
func (svc *PinpointEmail) Errorf(format string, v ...interface{}) {
	svc.logger.Errorf(serviceName, format, v...)
}
