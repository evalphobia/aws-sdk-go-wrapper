package mobileanalytics

import (
	SDK "github.com/aws/aws-sdk-go/service/mobileanalytics"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
	"github.com/evalphobia/aws-sdk-go-wrapper/private/errors"
)

const (
	serviceName = "MobileAnalytics"
)

// MobileAnalytics has MobileAnalytics client.
type MobileAnalytics struct {
	client *SDK.MobileAnalytics

	logger log.Logger
	prefix string
}

// New returns initialized *MobileAnalytics.
func New(conf config.Config) (*MobileAnalytics, error) {
	sess, err := conf.Session()
	if err != nil {
		return nil, err
	}

	svc := &MobileAnalytics{
		client: SDK.New(sess),
		logger: log.DefaultLogger,
		prefix: conf.DefaultPrefix,
	}
	return svc, nil
}

// SetLogger sets logger.
func (svc *MobileAnalytics) SetLogger(logger log.Logger) {
	svc.logger = logger
}

// PutRecord puts the given data into stream record.
func (svc *MobileAnalytics) PutEvents(list EventList) error {
	_, err := s.client.PutEvents(&SDK.PutEventsInput{
		Events: list.ToAWSEvents(),
	})
	if err != nil {
		s.service.Errorf("error on `PutEvents` operation; error=%s;", err.Error())
	}
	return err
}

// Infof logging information.
func (svc *MobileAnalytics) Infof(format string, v ...interface{}) {
	svc.logger.Infof(serviceName, format, v...)
}

// Errorf logging error information.
func (svc *MobileAnalytics) Errorf(format string, v ...interface{}) {
	svc.logger.Errorf(serviceName, format, v...)
}

func newErrors() *errors.Errors {
	return errors.NewErrors(serviceName)
}
