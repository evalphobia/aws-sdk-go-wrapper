// DynamoDB utility

package dynamodb

import (
	SDK "github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

// NewKeySchema creates new []*SDK.KeySchemaElement.
func NewKeySchema(elements ...*SDK.KeySchemaElement) []*SDK.KeySchemaElement {
	if len(elements) > 1 {
		schema := make([]*SDK.KeySchemaElement, 2)
		schema[0] = elements[0]
		schema[1] = elements[1]
		return schema
	}

	schema := make([]*SDK.KeySchemaElement, 1)
	schema[0] = elements[0]
	return schema
}

// NewKeyElement creates initialized *SDK.KeySchemaElement.
func NewKeyElement(keyName, keyType string) *SDK.KeySchemaElement {
	return &SDK.KeySchemaElement{
		AttributeName: pointers.String(keyName),
		KeyType:       pointers.String(keyType),
	}
}

// NewHashKeyElement creates initialized *SDK.KeySchemaElement for HashKey.
func NewHashKeyElement(keyName string) *SDK.KeySchemaElement {
	return NewKeyElement(keyName, KeyTypeHash)
}

// NewRangeKeyElement creates initialized *SDK.KeySchemaElement for RangeKey.
func NewRangeKeyElement(keyName string) *SDK.KeySchemaElement {
	return NewKeyElement(keyName, KeyTypeRange)
}
