// SNS App

package sns

import (
	"regexp"

	SDK "github.com/aws/aws-sdk-go/service/sns"

	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

type SNSApp struct {
	svc      *AmazonSNS
	platform string
	arn      string

	userData string
}

// SetUserData sets CustomUserData
func (a *SNSApp) SetUserData(userData string) {
	a.userData = userData
}

// Create Endpoint(add device) and return `EndpointARN`
func (a *SNSApp) createEndpoint(token string) (string, error) {
	in := &SDK.CreatePlatformEndpointInput{
		PlatformApplicationArn: String(a.arn),
		Token: String(token),
	}
	if a.userData != "" {
		in.CustomUserData = String(a.userData)
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

var reARNError = regexp.MustCompile("Endpoint (arn:aws:sns:.*) already exists")

func ParseARNFromError(err error) (arn string, ok bool) {
	if err == nil {
		return "", false
	}

	list := reARNError.FindStringSubmatch(err.Error())
	if len(list) == 2 {
		return list[1], true
	}
	return "", false
}
