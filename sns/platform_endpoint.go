// SNS endpoint

package sns

import (
	"strconv"

	SDK "github.com/aws/aws-sdk-go/service/sns"
)

// PlatformEndpoint is struct for Platform Endpoint.
type PlatformEndpoint struct {
	svc      *SNS
	arn      string
	protocol string
	token    string
	enable   bool
}

// Publish sends push notification to the endpoint.
func (e *PlatformEndpoint) Publish(msg string, badge int) error {
	return e.svc.Publish(e.arn, msg, map[string]interface{}{"badge": badge})
}

// PublishWithOption sends push notification to the endpoint with optional params.
func (e *PlatformEndpoint) PublishWithOption(msg string, opt map[string]interface{}) error {
	return e.svc.Publish(e.arn, msg, opt)
}

// GetARN returns endpoint ARN.
func (e *PlatformEndpoint) GetARN() (endpointARN string) {
	return e.arn
}

// GetToken returns endpoint Token.
func (e *PlatformEndpoint) GetToken() (token string) {
	return e.token
}

// Enable returns info that endpoint is Enable or not.
func (e *PlatformEndpoint) Enable() bool {
	return e.enable
}

// UpdateToken updates token and `Enabled` is true.
func (e *PlatformEndpoint) UpdateToken(token string) error {
	e.enable = true
	return e.updateToken(token, true)
}

// UpdateAsDisable updates as `Enabled` is false.
func (e *PlatformEndpoint) UpdateAsDisable() error {
	e.enable = false
	return e.updateToken(e.token, false)
}

// updateToken updates endpoint attributes.
func (e *PlatformEndpoint) updateToken(token string, isEnable bool) error {
	in := &SDK.SetEndpointAttributesInput{
		EndpointArn: String(e.arn),
		Attributes: map[string]*string{
			"Enabled": String(strconv.FormatBool(isEnable)),
			"Token":   String(token),
		},
	}
	_, err := e.svc.client.SetEndpointAttributes(in)
	if err != nil {
		e.svc.Errorf("error on `SetEndpointAttributes` operation; arn=%s; token=%s; error=%s;", e.arn, token, err.Error())
	}
	return err
}
