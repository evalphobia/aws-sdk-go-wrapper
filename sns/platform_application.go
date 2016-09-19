// SNS App

package sns

import SDK "github.com/aws/aws-sdk-go/service/sns"

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
	arn, err := a.createEndpoint(token, String(userData))
	if err != nil {
		return nil, err
	}
	endpoint := a.svc.newApplicationEndpoint(arn)
	return endpoint, nil
}

// createEndpoint creates Endpoint(add device) and return `EndpointARN`.
func (a *PlatformApplication) createEndpoint(token string, userData *string) (endpointARN string, err error) {
	in := &SDK.CreatePlatformEndpointInput{
		PlatformApplicationArn: String(a.arn),
		Token:          String(token),
		CustomUserData: userData,
	}

	resp, err := a.svc.client.CreatePlatformEndpoint(in)
	if err != nil {
		a.svc.Errorf("error on `CreatePlatformEndpoint` operation; arn=%s; token=%s; error=%s;", a.arn, token, err.Error())
		return "", err
	}
	return *resp.EndpointArn, nil
}
