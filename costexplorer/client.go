package costexplorer

import (
	"github.com/aws/aws-sdk-go/aws/session"
	SDK "github.com/aws/aws-sdk-go/service/costexplorer"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

const (
	serviceName = "CostExplorer"
)

// CostExplorer has *SDK.CostExplorer client.
type CostExplorer struct {
	client *SDK.CostExplorer

	logger log.Logger
}

// New returns initialized *CostExplorer.
func New(conf config.Config) (*CostExplorer, error) {
	sess, err := conf.Session()
	if err != nil {
		return nil, err
	}

	return NewFromSession(sess), nil
}

// NewFromSession returns initialized *CostExplorer from aws.Session.
func NewFromSession(sess *session.Session) *CostExplorer {
	return &CostExplorer{
		client: SDK.New(sess),
		logger: log.DefaultLogger,
	}
}

// SetLogger sets logger.
func (svc *CostExplorer) SetLogger(logger log.Logger) {
	svc.logger = logger
}

// GetCostAndUsage executes GetCostAndUsage operation with customized input.
func (svc *CostExplorer) GetCostAndUsage(input GetCostAndUsageInput) (UsageResult, error) {
	return svc.DoGetCostAndUsage(input.ToInput())
}

// DoGetCostAndUsage executes GetCostAndUsage operation.
func (svc *CostExplorer) DoGetCostAndUsage(input *SDK.GetCostAndUsageInput) (UsageResult, error) {
	output, err := svc.client.GetCostAndUsage(input)
	if err != nil {
		svc.Errorf("error on `GetCostAndUsage` operation; error=%w;", err)
		return UsageResult{}, err
	}

	return NewUsageResult(output), nil
}

// Infof logging information.
func (svc *CostExplorer) Infof(format string, v ...interface{}) {
	svc.logger.Infof(serviceName, format, v...)
}

// Errorf logging error information.
func (svc *CostExplorer) Errorf(format string, v ...interface{}) {
	svc.logger.Errorf(serviceName, format, v...)
}
