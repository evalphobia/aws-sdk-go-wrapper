package iam

import (
	"time"

	SDK "github.com/aws/aws-sdk-go/service/iam"
)

// Policy contains IAM policy data.
type Policy struct {
	ARN             string
	PolicyID        string
	PolicyName      string
	VersionID       string
	Description     string
	AttachmentCount int64
	CreateDate      time.Time
	UpdateDate      time.Time
}

// NewPoilicy returns initialized Policy from *SDK.Policy.
func NewPoilicy(p *SDK.Policy) Policy {
	pp := Policy{}
	if p.Arn != nil {
		pp.ARN = *p.Arn
	}
	if p.PolicyId != nil {
		pp.PolicyID = *p.PolicyId
	}
	if p.PolicyName != nil {
		pp.PolicyName = *p.PolicyName
	}
	if p.DefaultVersionId != nil {
		pp.VersionID = *p.DefaultVersionId
	}
	if p.Description != nil {
		pp.Description = *p.Description
	}
	if p.AttachmentCount != nil {
		pp.AttachmentCount = *p.AttachmentCount
	}
	if p.CreateDate != nil {
		pp.CreateDate = *p.CreateDate
	}
	if p.UpdateDate != nil {
		pp.UpdateDate = *p.UpdateDate
	}
	return pp
}

// NewPolicies converts from []*SDK.Policy to []Policy.
func NewPolicies(list []*SDK.Policy) []Policy {
	result := make([]Policy, len(list))
	for i, p := range list {
		result[i] = NewPoilicy(p)
	}
	return result
}
