// SNS topic

package sns

// Platforms contations Platforms Application ARNs.
type Platforms struct {
	Production bool

	Apple  string
	Google string
}

// GetARNByType returns ARN for application platform by device type.
func (p Platforms) GetARNByType(typ string) (arn string) {
	switch typ {
	case AppTypeAPNS, AppTypeAPNSSandbox:
		return p.Apple
	case AppTypeGCM:
		return p.Google
	}
	return ""
}
