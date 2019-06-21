package iam

import (
	SDK "github.com/aws/aws-sdk-go/service/iam"
)

// PolicyEntity contains Entity's id and name.
// Entity is User, Group or Role.
type PolicyEntity struct {
	Type EntityType
	ID   string
	Name string
}

// NewPolicyEntityList creates []PolicyEntity from *SDK.ListEntitiesForPolicyOutput.
func NewPolicyEntityList(o *SDK.ListEntitiesForPolicyOutput) []PolicyEntity {
	maxSize := len(o.PolicyUsers) + len(o.PolicyGroups) + len(o.PolicyRoles)
	list := make([]PolicyEntity, 0, maxSize)

	for _, e := range o.PolicyUsers {
		list = append(list, PolicyEntity{
			Type: NewEntityTypeUser(),
			ID:   *e.UserId,
			Name: *e.UserName,
		})
	}

	for _, e := range o.PolicyGroups {
		list = append(list, PolicyEntity{
			Type: NewEntityTypeGroup(),
			ID:   *e.GroupId,
			Name: *e.GroupName,
		})
	}

	for _, e := range o.PolicyRoles {
		list = append(list, PolicyEntity{
			Type: NewEntityTypeRole(),
			ID:   *e.RoleId,
			Name: *e.RoleName,
		})
	}
	return list
}

// IsUser checks this entity is user or not.
func (e PolicyEntity) IsUser() bool {
	return e.Type == entityTypeUser
}

// IsGroup checks this entity is group or not.
func (e PolicyEntity) IsGroup() bool {
	return e.Type == entityTypeGroup
}

// IsRole checks this entity is role or not.
func (e PolicyEntity) IsRole() bool {
	return e.Type == entityTypeRole
}

// EntityType represents entity's type.
type EntityType string

const (
	entityTypeUser  = "user"
	entityTypeGroup = "group"
	entityTypeRole  = "role"
)

// NewEntityTypeUser returns user's EntityType.
func NewEntityTypeUser() EntityType {
	return EntityType(entityTypeUser)
}

// NewEntityTypeGroup returns group's EntityType.
func NewEntityTypeGroup() EntityType {
	return EntityType(entityTypeGroup)
}

// NewEntityTypeRole returns role's EntityType.
func NewEntityTypeRole() EntityType {
	return EntityType(entityTypeRole)
}
