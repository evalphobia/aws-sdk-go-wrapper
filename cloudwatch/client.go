package cloudwatch

import (
	"github.com/aws/aws-sdk-go/aws/session"
	SDK "github.com/aws/aws-sdk-go/service/cloudwatch"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

const (
	serviceName = "CloudWatch"
)

// CloudWatch has CloudWatch client.
type CloudWatch struct {
	client *SDK.CloudWatch

	logger log.Logger
}

// New returns initialized *CloudWatch.
func New(conf config.Config) (*CloudWatch, error) {
	sess, err := conf.Session()
	if err != nil {
		return nil, err
	}

	return NewFromSession(sess), nil
}

// NewFromSession returns initialized *CloudWatch from aws.Session.
func NewFromSession(sess *session.Session) *CloudWatch {
	return &CloudWatch{
		client: SDK.New(sess),
		logger: log.DefaultLogger,
	}
}

// SetLogger sets logger.
func (svc *CloudWatch) SetLogger(logger log.Logger) {
	svc.logger = logger
}

// GetMetricStatistics executes GetMetricStatistics operation.
func (svc *CloudWatch) GetMetricStatistics(in MetricStatisticsInput) (*MetricStatisticsResponse, error) {
	out, err := svc.DoGetMetricStatistics(in.ToInput())
	if err != nil {
		return nil, err
	}
	return NewMetricStatisticsResponse(out), nil
}

// DoGetMetricStatistics executes GetMetricStatistics operation.
func (svc *CloudWatch) DoGetMetricStatistics(in *SDK.GetMetricStatisticsInput) (*SDK.GetMetricStatisticsOutput, error) {
	out, err := svc.client.GetMetricStatistics(in)
	if err != nil {
		svc.Errorf("error on `GetMetricStatistics` operation; error=%s;", err.Error())
		return nil, err
	}
	return out, nil
}

// PutMetricData executes PutMetricData operation.
func (svc *CloudWatch) PutMetricData(in PutMetricDataInput) error {
	return svc.DoPutMetricData(in.ToInput())
}

// DoPutMetricData executes PutMetricData operation.
func (svc *CloudWatch) DoPutMetricData(in *SDK.PutMetricDataInput) error {
	_, err := svc.client.PutMetricData(in)
	if err != nil {
		svc.Errorf("error on `PutMetricData` operation; error=%s;", err.Error())
	}
	return err
}

// Infof logging information.
func (svc *CloudWatch) Infof(format string, v ...interface{}) {
	svc.logger.Infof(serviceName, format, v...)
}

// Errorf logging error information.
func (svc *CloudWatch) Errorf(format string, v ...interface{}) {
	svc.logger.Errorf(serviceName, format, v...)
}
