// SNS client

package sns

import (
	AWS "github.com/awslabs/aws-sdk-go/aws"
	SNS "github.com/awslabs/aws-sdk-go/gen/sns"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

type SNSApp struct {
	client   *AmazonSNS
	platform string
	arn      string
}

// Create Endpoint(add device) and return `EndpointARN`
func (a *SNSApp) createEndpoint(token string) (string, error) {
	in := &SNS.CreatePlatformEndpointInput{
		PlatformApplicationARN: AWS.String(a.arn),
		Token: AWS.String(token),
	}
	resp, err := a.client.Client.CreatePlatformEndpoint(in)
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
	endpoint := &SNSEndpoint{
		arn:      arn,
		protocol: "application",
		client:   a.client,
	}
	return endpoint, nil
}
