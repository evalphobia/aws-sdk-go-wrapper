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

// Create Endpoint(add device) and return `EndpointARN`
func (a *SNSApp) createEndpoint(token string) (string, error) {
	in := &SDK.CreatePlatformEndpointInput{
		PlatformApplicationArn: String(a.arn),
		Token: String(token),
	}
	resp, err := a.svc.Client.CreatePlatformEndpoint(in)
	if err != nil {
		log.Error("[SNS] error on `CreatePlatformEndpoint` operation, token="+token, err.Error())
		return "", err
	}
	return *resp.EndpointArn, nil
}

// Create Endpoint(add device) and return `EndpointARN`
func (a *SNSApp) CreateEndpoint(token string) (*SNSEndpoint, error) {
	arn, err := a.createEndpoint(token)
	if err != nil {
		return nil, err
	}
	endpoint := a.svc.NewApplicationEndpoint(arn)
	return endpoint, nil
}
