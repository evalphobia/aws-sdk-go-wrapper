// SNS endpoint

package sns

type SNSEndpoint struct {
	arn      string
	protocol string
	svc      *AmazonSNS
}

func NewEndpoint(arn, protocol string, svc *AmazonSNS) *SNSEndpoint {
	return &SNSEndpoint{
		arn:      arn,
		protocol: protocol,
		svc:      svc,
	}
}

// Publish notification to the endpoint
func (e *SNSEndpoint) Publish(msg string, badge int) error {
	return e.svc.Publish(e.arn, msg, map[string]interface{}{"badge": badge})
}

// return Endpoint ARN
func (e *SNSEndpoint) GetARN() string {
	return e.arn
}
