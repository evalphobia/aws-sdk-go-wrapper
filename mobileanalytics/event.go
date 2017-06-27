package mobileanalytics

import (
	"time"

	SDK "github.com/aws/aws-sdk-go/service/mobileanalytics"
)

type Event struct {
	Attributes map[string]string
	EventType  string
	Metrics    map[string]float64
	Session interface{}
    Timestamp　time.Time
	Version string
}

func (e *Event) ToAWSEvent() *SDK.Event {
	event := &SDK.Event{}
	if e.EventType != "" {
		event.SetEventType(e.EventType)
	}
	if e.Version != "" {
		event.SetVersion(e.Version)
	}
	if !e.Timestamp.IsZero() {
		event.SetTimestamp(e.Timestamp.String())
	}
	if e.Session != nil {
		// event.SetSession(e.Session)
	}



	if len(e.Attributes) != 0 {
		attr = make(map[string]*string)
		for k, v := range e.Attributes {
			attr[k] = &v
		}
		event.SetAttributes(attr)
	}

	if len(e.Metrics) != 0 {
		metrics = make(map[string]*float64)
		for k, v := range e.Metrics {
			metrics[k] = &v
		}
		event.SetMetrics(metrics)
	}

	return event
}



type EventList struct {
	List []*Event
}

func (l EventList) ToAWSEvents() []*SDK.Event {
	events := make([]*SDK.Event, len(l.List))
	for i, event := range l.List {
		events[i] = event.ToAWSEvent()
	}
	return events
}




