// DynamoDB utility

package dynamodb

import (
	SDK "github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

type AttributeDefinition struct {
	Name string
	Type string
}

func NewAttributeDefinition(out *SDK.AttributeDefinition) AttributeDefinition {
	v := AttributeDefinition{}
	if out == nil {
		return v
	}

	if out.AttributeName != nil {
		v.Name = *out.AttributeName
	}
	if out.AttributeType != nil {
		v.Type = *out.AttributeType
	}
	return v
}

func NewAttributeDefinitionFromType(name, typ string) AttributeDefinition {
	return AttributeDefinition{
		Name: name,
		Type: typ,
	}
}

func (d AttributeDefinition) IsEmpty() bool {
	switch {
	case d.Name != "",
		d.Type != "":
		return false
	}
	return true
}

func (d AttributeDefinition) ToSDKType() *SDK.AttributeDefinition {
	if d.IsEmpty() {
		return nil
	}
	return &SDK.AttributeDefinition{
		AttributeName: pointers.String(d.Name),
		AttributeType: pointers.String(d.Type),
	}
}

func NewAttributeDefinitionList(list []*SDK.AttributeDefinition) []AttributeDefinition {
	if len(list) == 0 {
		return nil
	}

	result := make([]AttributeDefinition, len(list))
	for i, out := range list {
		result[i] = NewAttributeDefinition(out)
	}
	return result
}

// NewStringAttribute returns a table AttributeDefinition for string.
func NewStringAttribute(attrName string) AttributeDefinition {
	return NewAttributeDefinitionFromType(attrName, AttributeTypeString)
}

// NewNumberAttribute returns a table AttributeDefinition for number.
func NewNumberAttribute(attrName string) AttributeDefinition {
	return NewAttributeDefinitionFromType(attrName, AttributeTypeNumber)
}

// NewByteAttribute returns a table AttributeDefinition for byte.
func NewByteAttribute(attrName string) AttributeDefinition {
	return NewAttributeDefinitionFromType(attrName, AttributeTypeBinary)
}

// NewBoolAttribute returns a table AttributeDefinition for boolean.
func NewBoolAttribute(attrName string) AttributeDefinition {
	return NewAttributeDefinitionFromType(attrName, AttributeTypeBool)
}

// NewNullAttribute returns a table AttributeDefinition for null.
func NewNullAttribute(attrName string) AttributeDefinition {
	return NewAttributeDefinitionFromType(attrName, AttributeTypeNull)
}

// NewMapAttribute returns a table AttributeDefinition for map.
func NewMapAttribute(attrName string) AttributeDefinition {
	return NewAttributeDefinitionFromType(attrName, AttributeTypeMap)
}

// NewListAttribute returns a table AttributeDefinition for list.
func NewListAttribute(attrName string) AttributeDefinition {
	return NewAttributeDefinitionFromType(attrName, AttributeTypeList)
}

// NewStringSetAttribute returns a table AttributeDefinition for string set.
func NewStringSetAttribute(attrName string) AttributeDefinition {
	return NewAttributeDefinitionFromType(attrName, AttributeTypeStringSet)
}

// NewNumberSetAttribute returns a table AttributeDefinition for number set.
func NewNumberSetAttribute(attrName string) AttributeDefinition {
	return NewAttributeDefinitionFromType(attrName, AttributeTypeNumberSet)
}

// NewBinarySetAttribute returns a table AttributeDefinition for binary set.
func NewBinarySetAttribute(attrName string) AttributeDefinition {
	return NewAttributeDefinitionFromType(attrName, AttributeTypeBinarySet)
}
