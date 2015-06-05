// SNS App

package sns

import (
	SDK "github.com/aws/aws-sdk-go/service/sns"

	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

type SNSApp struct {
	svc      *AmazonSNS
	platform string
	arn      string
}

func NewApp(arn, pf string, svc *AmazonSNS) *SNSApp {
	return &SNSApp{
		arn:      arn,
		platform: pf,
		svc:      svc,
	}
}

// Create Endpoint(add device) and return `EndpointARN`
func (a *SNSApp) createEndpoint(token string) (string, error) {
	in := &SDK.CreatePlatformEndpointInput{
		PlatformApplicationARN: String(a.arn),
		Token: String(token),
	}
	resp, err := a.svc.Client.CreatePlatformEndpoint(in)
	if err != nil {
		log.Error("[SNS] error on `CreatePlatformEndpoint` operation, token="+token, err.Error())
		return "", err
	}
	return *resp.EndpointARN, nil
}

// Create Endpoint(add device) and return `EndpointARN`
func (a *SNSApp) CreateEndpoint(token string) (*SNSEndpoint, error) {
	arn, err := a.createEndpoint(token)
	if err != nil {
		return nil, err
	}
	endpoint := NewEndpoint(arn, "application", a.svc)
	return endpoint, nil
}
