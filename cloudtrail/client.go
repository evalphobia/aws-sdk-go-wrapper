package cloudtrail

import (
	SDK "github.com/aws/aws-sdk-go/service/cloudtrail"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

const (
	serviceName = "CloudTrail"
)

// CloudTrail has CloudTrail client.
type CloudTrail struct {
	client *SDK.CloudTrail

	logger log.Logger
}

// New returns initialized *CloudTrail.
func New(conf config.Config) (*CloudTrail, error) {
	sess, err := conf.Session()
	if err != nil {
		return nil, err
	}

	svc := &CloudTrail{
		client: SDK.New(sess),
		logger: log.DefaultLogger,
	}
	return svc, nil
}

// SetLogger sets logger.
func (svc *CloudTrail) SetLogger(logger log.Logger) {
	svc.logger = logger
}

// LookupEventsAll executes LookupEvents operation and fetch all of events.
func (svc *CloudTrail) LookupEventsAll(input LookupEventsInput) (LookupEventsResult, error) {
	in := input
	result := LookupEventsResult{}
	for {
		res, err := svc.LookupEvents(in)
		if err != nil {
			return result, err
		}

		result.Events = append(result.Events, res.Events...)
		if res.NextToken == "" {
			return result, nil
		}
		in.NextToken = res.NextToken
	}
}

// LookupEvents executes LookupEvents operation with customized input.
func (svc *CloudTrail) LookupEvents(input LookupEventsInput) (LookupEventsResult, error) {
	return svc.DoLookupEvents(input.ToInput())
}

// DoLookupEvents executes LookupEvents operation.
func (svc *CloudTrail) DoLookupEvents(in *SDK.LookupEventsInput) (LookupEventsResult, error) {
	out, err := svc.client.LookupEvents(in)
	if err != nil {
		svc.Errorf("error on `LookupEvents` operation; error=%s;", err.Error())
		return LookupEventsResult{}, err
	}
	return NewLookupEventsResult(out), nil
}

// Infof logging information.
func (svc *CloudTrail) Infof(format string, v ...interface{}) {
	svc.logger.Infof(serviceName, format, v...)
}

// Errorf logging error information.
func (svc *CloudTrail) Errorf(format string, v ...interface{}) {
	svc.logger.Errorf(serviceName, format, v...)
}
