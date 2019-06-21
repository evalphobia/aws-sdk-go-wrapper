package iam

import (
	"time"

	SDK "github.com/aws/aws-sdk-go/service/iam"
)

// Group contains IAM Group data.
type Group struct {
	ARN        string
	GroupID    string
	GroupName  string
	Path       string
	CreateDate time.Time
}

// NewGroup returns initialized Group from *SDK.Group.
func NewGroup(g *SDK.Group) Group {
	gg := Group{}
	if g.Arn != nil {
		gg.ARN = *g.Arn
	}
	if g.GroupId != nil {
		gg.GroupID = *g.GroupId
	}
	if g.GroupName != nil {
		gg.GroupName = *g.GroupName
	}
	if g.Path != nil {
		gg.Path = *g.Path
	}
	if g.CreateDate != nil {
		gg.CreateDate = *g.CreateDate
	}
	return gg
}

// NewGroups converts from []*SDK.Group to []Group.
func NewGroups(list []*SDK.Group) []Group {
	result := make([]Group, len(list))
	for i, p := range list {
		result[i] = NewGroup(p)
	}
	return result
}
