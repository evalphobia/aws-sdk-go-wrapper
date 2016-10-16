// DynamoDB utility

package dynamodb

import (
	SDK "github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

// NewAttributeDefinitions returns multiple AttributeDefinition to single slice.
func NewAttributeDefinitions(attr ...*SDK.AttributeDefinition) []*SDK.AttributeDefinition {
	return attr
}

// NewAttributeDefinition returns initialized *SDK.AttributeDefinition.
func NewAttributeDefinition(attrName, attrType string) *SDK.AttributeDefinition {
	newAttr := &SDK.AttributeDefinition{}
	var typ *string
	switch attrType {
	case "S", "N", "B", "BOOL", "L", "M", "SS", "NS", "BS":
		typ = pointers.String(attrType)
	default:
		return newAttr
	}
	newAttr.AttributeName = pointers.String(attrName)
	newAttr.AttributeType = typ
	return newAttr
}

// NewStringAttribute returns a table AttributeDefinition for string
func NewStringAttribute(attrName string) *SDK.AttributeDefinition {
	return NewAttributeDefinition(attrName, "S")
}

// NewNumberAttribute returns a table AttributeDefinition for number
func NewNumberAttribute(attrName string) *SDK.AttributeDefinition {
	return NewAttributeDefinition(attrName, "N")
}

// NewByteAttribute returns a table AttributeDefinition for byte
func NewByteAttribute(attrName string) *SDK.AttributeDefinition {
	return NewAttributeDefinition(attrName, "B")
}

// NewBoolAttribute returns a table AttributeDefinition for boolean
func NewBoolAttribute(attrName string) *SDK.AttributeDefinition {
	return NewAttributeDefinition(attrName, "BOOL")
}
