// SNS endpoint

package sns

import (
	"strconv"

	SDK "github.com/aws/aws-sdk-go/service/sns"
)

// SNSEndpoint is struct for endpoint(device / )
type SNSEndpoint struct {
	svc      *AmazonSNS
	arn      string
	protocol string
	token    string
	enable   bool
}

// Publish sends push notification to the endpoint
func (e *SNSEndpoint) Publish(msg string, badge int) error {
	return e.svc.Publish(e.arn, msg, map[string]interface{}{"badge": badge})
}

// PublishWithOption sends push notification to the endpoint with optional params
func (e *SNSEndpoint) PublishWithOption(msg string, opt map[string]interface{}) error {
	return e.svc.Publish(e.arn, msg, opt)
}

// GetARN returns endpoint ARN
func (e *SNSEndpoint) GetARN() string {
	return e.arn
}

// GetToken returns endpoint Token
func (e *SNSEndpoint) GetToken() string {
	return e.token
}

// SetToken set endpoint Token
func (e *SNSEndpoint) SetToken(token string) {
	e.token = token
}

// Enable returns endpoint Enable
func (e *SNSEndpoint) Enable() bool {
	return e.enable
}

// UpdateTokenAsEnable updates token and enabled as true
func (e *SNSEndpoint) UpdateTokenAsEnable() error {
	e.enable = true
	return e.UpdateToken()
}

// UpdateTokenAsEnable updates token and enabled as false
func (e *SNSEndpoint) UpdateTokenAsDisable() error {
	e.enable = false
	return e.UpdateToken()
}

// UpdateToken updates endpoint attributes
func (e *SNSEndpoint) UpdateToken() error {
	in := &SDK.SetEndpointAttributesInput{
		EndpointArn: String(e.arn),
		Attributes: map[string]*string{
			"Enabled": String(strconv.FormatBool(e.enable)),
			"Token":   String(e.token),
		},
	}
	_, err := e.svc.Client.SetEndpointAttributes(in)
	return err
}
