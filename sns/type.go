package sns

import (
	"strconv"
	"time"
)

// PlatformAttributes contains platform application attributes.
type PlatformAttributes struct {
	HasPlatformCredential bool
	PlatformCredential    string

	HasPlatformPrincipal bool
	PlatformPrincipal    string

	HasEventEndpointCreated bool
	EventEndpointCreated    string

	HasEventEndpointDeleted bool
	EventEndpointDeleted    string

	HasEventEndpointUpdated bool
	EventEndpointUpdated    string

	HasEventDeliveryFailure bool
	EventDeliveryFailure    string

	HasSuccessFeedbackRoleArn bool
	SuccessFeedbackRoleArn    string

	HasFailureFeedbackRoleArn bool
	FailureFeedbackRoleArn    string

	HasSuccessFeedbackSampleRate bool
	SuccessFeedbackSampleRate    int // 0 - 100

	HasEnabled bool
	Enabled    bool

	HasAppleCertificateExpirationDate bool
	AppleCertificateExpirationDate    time.Time
}

func NewPlatformAttributesFromMap(attr map[string]*string) PlatformAttributes {
	a := PlatformAttributes{}
	if len(attr) == 0 {
		return a
	}

	if v, ok := attr["PlatformCredential"]; ok {
		a.HasPlatformCredential = true
		if v != nil {
			a.PlatformCredential = *v
		}
	}
	if v, ok := attr["PlatformPrincipal"]; ok {
		a.HasPlatformPrincipal = true
		if v != nil {
			a.PlatformPrincipal = *v
		}
	}
	if v, ok := attr["EventEndpointCreated"]; ok {
		a.HasEventEndpointCreated = true
		if v != nil {
			a.EventEndpointCreated = *v
		}
	}
	if v, ok := attr["EventEndpointDeleted"]; ok {
		a.HasEventEndpointDeleted = true
		if v != nil {
			a.EventEndpointDeleted = *v
		}
	}
	if v, ok := attr["EventEndpointUpdated"]; ok {
		a.HasEventEndpointUpdated = true
		if v != nil {
			a.EventEndpointUpdated = *v
		}
	}
	if v, ok := attr["EventDeliveryFailure"]; ok {
		a.HasEventDeliveryFailure = true
		if v != nil {
			a.EventDeliveryFailure = *v
		}
	}
	if v, ok := attr["SuccessFeedbackRoleArn"]; ok {
		a.HasSuccessFeedbackRoleArn = true
		if v != nil {
			a.SuccessFeedbackRoleArn = *v
		}
	}
	if v, ok := attr["FailureFeedbackRoleArn"]; ok {
		a.HasFailureFeedbackRoleArn = true
		if v != nil {
			a.FailureFeedbackRoleArn = *v
		}
	}
	if v, ok := attr["SuccessFeedbackSampleRate"]; ok {
		a.HasSuccessFeedbackSampleRate = true
		if v != nil {
			num, _ := strconv.Atoi(*v)
			a.SuccessFeedbackSampleRate = num
		}
	}
	if v, ok := attr["Enabled"]; ok {
		a.HasEnabled = true
		if v != nil {
			ok, _ := strconv.ParseBool(*v)
			a.Enabled = ok
		}
	}
	if v, ok := attr["AppleCertificateExpirationDate"]; ok {
		a.HasAppleCertificateExpirationDate = true
		if v != nil {
			dt, _ := time.Parse(time.RFC3339, *v)
			a.AppleCertificateExpirationDate = dt
		}
	}
	return a
}
