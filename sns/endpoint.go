// SNS client

package sns

type SNSEndpoint struct {
	arn      string
	protocol string
	client   *AmazonSNS
}

// Publish notification to the endpoint
func (e *SNSEndpoint) Publish(msg string, badge int) error {
	return e.client.Publish(e.arn, msg, map[string]interface{}{"badge": badge})
}
