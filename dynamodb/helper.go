// DynamoDB utility

package dynamodb

import (
	AWS "github.com/awslabs/aws-sdk-go/aws"
	DynamoDB "github.com/awslabs/aws-sdk-go/gen/dynamodb"

	"fmt"
	"strconv"
)

type Any interface{}

// Create new AttributeValue from the type of value
func createAttributeValue(v Any) DynamoDB.AttributeValue {
	switch t := v.(type) {
	case string:
		return DynamoDB.AttributeValue{
			S: AWS.String(t),
		}
	case int, int32, int64, uint, uint32, uint64, float32, float64:
		return DynamoDB.AttributeValue{
			N: AWS.String(fmt.Sprint(t)),
		}
	default:
		return DynamoDB.AttributeValue{}
	}
}

// Retrieve value from DynamoDB type
func getItemValue(val DynamoDB.AttributeValue) Any {
	switch {
	case val.N != nil:
		data, _ := strconv.Atoi(*val.N)
		return data
	case val.S != nil:
		return *val.S
	case val.BOOL != nil:
		return *val.BOOL
	case len(val.B) > 0:
		return val.B
	case len(val.M) > 0:
		return Unmarshal(val.M)
	case len(val.NS) > 0:
		var data []int
		for _, vString := range val.NS {
			vInt, _ := strconv.Atoi(vString)
			data = append(data, vInt)
		}
		return data
	case len(val.SS) > 0:
		var data []string
		for _, vString := range val.SS {
			data = append(data, vString)
		}
		return data
	case len(val.BS) > 0:
		var data [][]byte
		for _, vBytes := range val.BS {
			data = append(data, vBytes)
		}
		return data
	case len(val.L) > 0:
		var data []interface{}
		for _, vAny := range val.L {
			data = append(data, getItemValue(vAny))
		}
		return data
	}
	return nil
}

// Convert DynamoDB Item to map data
func Unmarshal(item map[string]DynamoDB.AttributeValue) map[string]interface{} {
	data := make(map[string]interface{})
	for key, val := range item {
		data[key] = getItemValue(val)
	}
	return data
}

// Create new KeySchema slice
func NewKeySchema(elements ...*DynamoDB.KeySchemaElement) []DynamoDB.KeySchemaElement {
	if len(elements) > 1 {
		schema := make([]DynamoDB.KeySchemaElement, 2, 2)
		schema[0] = *elements[0]
		schema[1] = *elements[1]
		return schema
	} else {
		schema := make([]DynamoDB.KeySchemaElement, 1, 1)
		schema[0] = *elements[0]
		return schema
	}
}

// Create new single KeySchema
func NewKeyElement(keyName, keyType string) *DynamoDB.KeySchemaElement {
	return &DynamoDB.KeySchemaElement{
		AttributeName: AWS.String(keyName),
		KeyType:       AWS.String(keyType),
	}
}

// Create new single KeySchema for HashKey
func NewHashKeyElement(keyName string) *DynamoDB.KeySchemaElement {
	return NewKeyElement(keyName, DynamoDB.KeyTypeHash)
}

// Create new single KeySchema for RangeKey
func NewRangeKeyElement(keyName string) *DynamoDB.KeySchemaElement {
	return NewKeyElement(keyName, DynamoDB.KeyTypeRange)
}

// Convert multiple definition to single slice
func NewAttributeDefinitions(attr ...DynamoDB.AttributeDefinition) []DynamoDB.AttributeDefinition {
	return attr
}

// Create new definition of table
func NewAttributeDefinition(attrName, attrType string) DynamoDB.AttributeDefinition {
	newAttr := DynamoDB.AttributeDefinition{}
	var typ *string
	switch attrType {
	case "S", "N", "B", "BOOL", "L", "M", "SS", "NS", "BS":
		typ = AWS.String(attrType)
	default:
		return newAttr
	}
	newAttr.AttributeName = AWS.String(attrName)
	newAttr.AttributeType = typ
	return newAttr
}

// Create new definition of table for string
func NewStringAttribute(attrName string) DynamoDB.AttributeDefinition {
	return NewAttributeDefinition(attrName, "S")
}

// Create new definition of table for number
func NewNumberAttribute(attrName string) DynamoDB.AttributeDefinition {
	return NewAttributeDefinition(attrName, "N")
}

// Create new definition of table for byte
func NewByteAttribute(attrName string) DynamoDB.AttributeDefinition {
	return NewAttributeDefinition(attrName, "B")
}

// Create new definition of table for boolean
func NewBoolAttribute(attrName string) DynamoDB.AttributeDefinition {
	return NewAttributeDefinition(attrName, "BOOL")
}
