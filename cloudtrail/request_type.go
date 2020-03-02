package cloudtrail

import (
	"time"

	SDK "github.com/aws/aws-sdk-go/service/cloudtrail"
)

// LookupEventsInput is optional parameters for `LookupEvents`.
type LookupEventsInput struct {
	StartTime     time.Time
	EndTime       time.Time
	MaxResults    int64 // 1~50
	EventCategory string
	NextToken     string

	LookupAttributes []LookupAttribute
}

func (o LookupEventsInput) ToInput() *SDK.LookupEventsInput {
	in := &SDK.LookupEventsInput{}
	if !o.StartTime.IsZero() {
		in.StartTime = &o.StartTime
	}
	if !o.EndTime.IsZero() {
		in.EndTime = &o.EndTime
	}
	if o.MaxResults != 0 {
		in.MaxResults = &o.MaxResults
	}
	if o.EventCategory != "" {
		in.EventCategory = &o.EventCategory
	}
	if o.NextToken != "" {
		in.NextToken = &o.NextToken
	}

	for _, v := range o.LookupAttributes {
		in.LookupAttributes = append(in.LookupAttributes, v.ToSDKValue())
	}
	return in
}

type LookupAttribute struct {
	Key   string
	Value string
}

func (a LookupAttribute) ToSDKValue() *SDK.LookupAttribute {
	return &SDK.LookupAttribute{
		AttributeKey:   &a.Key,
		AttributeValue: &a.Value,
	}
}
