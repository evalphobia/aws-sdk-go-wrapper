package cloudtrail

import (
	"time"

	SDK "github.com/aws/aws-sdk-go/service/cloudtrail"
)

// LookupEventsResult represents results from `LookupEvents` operation.
type LookupEventsResult struct {
	Events    []Event
	NextToken string
}

func NewLookupEventsResult(output *SDK.LookupEventsOutput) LookupEventsResult {
	r := LookupEventsResult{}
	if output == nil {
		return r
	}

	if output.NextToken != nil {
		r.NextToken = *output.NextToken
	}

	r.Events = newEvents(output.Events)
	return r
}

type Event struct {
	AccessKeyID     string
	CloudTrailEvent string
	EventID         string
	EventName       string
	EventSource     string
	EventTime       time.Time
	ReadOnly        string
	Username        string

	Resources []Resource
}

func newEvents(list []*SDK.Event) []Event {
	if len(list) == 0 {
		return nil
	}

	results := make([]Event, len(list))
	for i, v := range list {
		results[i] = newEvent(v)
	}
	return results
}

func newEvent(d *SDK.Event) Event {
	result := Event{}
	if d == nil {
		return result
	}

	if d.AccessKeyId != nil {
		result.AccessKeyID = *d.AccessKeyId
	}
	if d.CloudTrailEvent != nil {
		result.CloudTrailEvent = *d.CloudTrailEvent
	}
	if d.EventId != nil {
		result.EventID = *d.EventId
	}
	if d.EventName != nil {
		result.EventName = *d.EventName
	}
	if d.EventSource != nil {
		result.EventSource = *d.EventSource
	}
	if d.EventTime != nil {
		result.EventTime = *d.EventTime
	}
	if d.ReadOnly != nil {
		result.ReadOnly = *d.ReadOnly
	}
	if d.Username != nil {
		result.Username = *d.Username
	}

	result.Resources = newResources(d.Resources)
	return result
}

type Resource struct {
	Name string
	Type string
}

func newResources(list []*SDK.Resource) []Resource {
	if len(list) == 0 {
		return nil
	}

	results := make([]Resource, len(list))
	for i, v := range list {
		results[i] = newResource(v)
	}
	return results
}

func newResource(d *SDK.Resource) Resource {
	result := Resource{}
	if d == nil {
		return result
	}

	if d.ResourceName != nil {
		result.Name = *d.ResourceName
	}
	if d.ResourceType != nil {
		result.Type = *d.ResourceType
	}
	return result
}
