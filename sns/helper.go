// SNS client

package sns

import (
	"fmt"
)

const (
	// use static templete, converting map data to json is slow
	messageTemplateGCM         = `{"data": {"message": "%s"}}`
	messageTemplateAPNS        = `{"aps":{"alert": "%s", "sound": "%s"}}`
	messageTemplateAPNSBadge   = `{"aps":{"alert": "%s", "sound": "%s", "badge": %d}}`
	messageTemplateAPNSSandbox = `{"aps": {"alert": "%s"}}`
)

// make sns message for Google Cloud Messaging
func composeMessageGCM(msg string) string {
	return fmt.Sprintf(messageTemplateGCM, msg)
}

// make sns message for Apple Push Notification Service
func composeMessageAPNS(msg string, opt map[string]interface{}) string {
	b, hasBadge := opt["badge"]
	s, hasSound := opt["sound"]

	switch {
	case hasBadge && hasSound:
		return fmt.Sprintf(messageTemplateAPNSBadge, msg, s, b)
	case hasBadge:
		return fmt.Sprintf(messageTemplateAPNSBadge, msg, "default", b)
	case hasSound:
		return fmt.Sprintf(messageTemplateAPNS, msg, s)
	default:
		return fmt.Sprintf(messageTemplateAPNS, msg, "default")
	}
}

// make sns message for Apple Push Notification Service Sandbox
func composeMessageAPNSSandbox(msg string) string {
	return fmt.Sprintf(messageTemplateAPNSSandbox, msg)
}
