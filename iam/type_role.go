package iam

import (
	"time"

	SDK "github.com/aws/aws-sdk-go/service/iam"
)

// Role contains IAM Role data.
type Role struct {
	ARN                      string
	RoleID                   string
	RoleName                 string
	Path                     string
	Description              string
	AssumeRolePolicyDocument string
	CreateDate               time.Time
}

// NewRole returns initialized Role from *SDK.Role.
func NewRole(r *SDK.Role) Role {
	rr := Role{}
	if r.Arn != nil {
		rr.ARN = *r.Arn
	}
	if r.RoleId != nil {
		rr.RoleID = *r.RoleId
	}
	if r.RoleName != nil {
		rr.RoleName = *r.RoleName
	}
	if r.Path != nil {
		rr.Path = *r.Path
	}
	if r.CreateDate != nil {
		rr.CreateDate = *r.CreateDate
	}
	return rr
}

// NewRoles converts from []*SDK.Role to []Role.
func NewRoles(list []*SDK.Role) []Role {
	result := make([]Role, len(list))
	for i, p := range list {
		result[i] = NewRole(p)
	}
	return result
}
