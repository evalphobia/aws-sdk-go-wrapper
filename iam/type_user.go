package iam

import (
	"time"

	SDK "github.com/aws/aws-sdk-go/service/iam"
)

// User contains IAM User data.
type User struct {
	ARN              string
	UserID           string
	UserName         string
	Path             string
	CreateDate       time.Time
	PasswordLastUsed time.Time
}

// NewUser returns initialized User from *SDK.User.
func NewUser(u *SDK.User) User {
	uu := User{}
	if u.Arn != nil {
		uu.ARN = *u.Arn
	}
	if u.UserId != nil {
		uu.UserID = *u.UserId
	}
	if u.UserName != nil {
		uu.UserName = *u.UserName
	}
	if u.Path != nil {
		uu.Path = *u.Path
	}
	if u.CreateDate != nil {
		uu.CreateDate = *u.CreateDate
	}
	if u.PasswordLastUsed != nil {
		uu.PasswordLastUsed = *u.PasswordLastUsed
	}
	return uu
}

// NewUsers converts from []*SDK.User to []User.
func NewUsers(list []*SDK.User) []User {
	result := make([]User, len(list))
	for i, p := range list {
		result[i] = NewUser(p)
	}
	return result
}
