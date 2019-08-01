package cloudwatch

import (
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

	svc := &CloudWatch{
		client: SDK.New(sess),
		logger: log.DefaultLogger,
	}
	return svc, nil
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

// Infof logging information.
func (svc *CloudWatch) Infof(format string, v ...interface{}) {
	svc.logger.Infof(serviceName, format, v...)
}

// Errorf logging error information.
func (svc *CloudWatch) Errorf(format string, v ...interface{}) {
	svc.logger.Errorf(serviceName, format, v...)
}
