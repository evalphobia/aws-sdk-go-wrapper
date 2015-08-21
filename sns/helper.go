// SNS client

package sns

import (
	"encoding/json"

	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

const (
	gcmKeyMessage  = "message"
	apnsKeyMessage = "alert"
	apnsKeySound   = "sound"
	apnsKeyBadge   = "badge"
)

// make sns message for Google Cloud Messaging
func composeMessageGCM(msg string, opt map[string]interface{}) string {
	data := make(map[string]interface{})
	data[gcmKeyMessage] = msg
	for k, v := range opt {
		data[k] = v
	}

	message := make(map[string]interface{})
	message["data"] = data

	payload, err := json.Marshal(message)
	if err != nil {
		log.Error("[SNS] error on json.Marshal", err.Error())
	}
	return string(payload)
}

// make sns message for Apple Push Notification Service
func composeMessageAPNS(msg string, opt map[string]interface{}) string {
	aps := make(map[string]interface{})
	aps[apnsKeyMessage] = msg

	aps[apnsKeySound] = "default"
	if v, ok := opt[apnsKeySound]; ok {
		aps[apnsKeySound] = v
	}

	if v, ok := opt[apnsKeyBadge]; ok {
		aps[apnsKeyBadge] = v
	}

	message := make(map[string]interface{})
	message["aps"] = aps
	for k, v := range opt {
		switch k {
		case apnsKeySound:
			continue
		case apnsKeyBadge:
			continue
		default:
			message[k] = v
		}
	}

	payload, err := json.Marshal(message)
	if err != nil {
		log.Error("[SNS] error on json.Marshal", err.Error())
	}
	return string(payload)
}
