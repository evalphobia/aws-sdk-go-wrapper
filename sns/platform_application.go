package sns

import (
	"regexp"

	SDK "github.com/aws/aws-sdk-go/service/sns"

	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

// PlatformApplication is struct for Platform Application.
type PlatformApplication struct {
	svc      *SNS
	platform string
	arn      string

	userData string
}

// SetUserData sets CustomUserData.
func (a *PlatformApplication) SetUserData(userData string) {
	a.userData = userData
}

// CreateEndpoint creates Endpoint(add device).
func (a *PlatformApplication) CreateEndpoint(token string) (*PlatformEndpoint, error) {
	arn, err := a.createEndpoint(token, nil)
	if err != nil {
		return nil, err
	}
	endpoint := a.svc.newApplicationEndpoint(arn)
	return endpoint, nil
}

// CreateEndpointWithUserData creates Endpoint(add device) with CustomUserData.
func (a *PlatformApplication) CreateEndpointWithUserData(token, userData string) (*PlatformEndpoint, error) {
	arn, err := a.createEndpoint(token, pointers.String(userData))
	if err != nil {
		return nil, err
	}
	endpoint := a.svc.newApplicationEndpoint(arn)
	return endpoint, nil
}

// createEndpoint creates Endpoint(add device) and return `EndpointARN`.
func (a *PlatformApplication) createEndpoint(token string, userData *string) (endpointARN string, err error) {
	in := &SDK.CreatePlatformEndpointInput{
		PlatformApplicationArn: pointers.String(a.arn),
		Token:                  pointers.String(token),
		CustomUserData:         userData,
	}

	resp, err := a.svc.client.CreatePlatformEndpoint(in)
	if err != nil {
		a.svc.Errorf("error on `CreatePlatformEndpoint` operation; arn=%s; token=%s; error=%s;", a.arn, token, err.Error())
		return "", err
	}
	return *resp.EndpointArn, nil
}

var reARNError = regexp.MustCompile("Endpoint (arn:aws:sns:.*) already exists")

// ParseARNFromError extracts ARN string from error message.
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
