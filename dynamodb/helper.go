// DynamoDB utility

package dynamodb

import (
	"fmt"
	"reflect"
	"strconv"

	SDK "github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	KeyTypeHash  = "HASH"
	KeyTypeRange = "RANGE"
)

type Any interface{}

func newAttributeValue(typ string, val interface{}) *SDK.AttributeValue {
	switch typ {
	case "S":
		return newAttributeValueS(val)
	case "N":
		return newAttributeValueN(val)
	case "B":
		return newAttributeValueB(val)
	case "BOOL":
		return newAttributeValueBOOL(val)
	case "SS":
		return newAttributeValueSS(val)
	case "NS":
		return newAttributeValueNS(val)
	case "L":
		return newAttributeValueL(val)
	case "M":
		return newAttributeValueM(val)
	}
	return nil
}

func newAttributeValueS(val interface{}) *SDK.AttributeValue {
	return &SDK.AttributeValue{S: String(fmt.Sprint(val))}
}

func newAttributeValueN(val interface{}) *SDK.AttributeValue {
	return &SDK.AttributeValue{N: String(fmt.Sprint(val))}
}

func newAttributeValueB(val interface{}) *SDK.AttributeValue {
	switch t := val.(type) {
	case []byte:
		return &SDK.AttributeValue{B: t}
	}
	return nil
}

func newAttributeValueBOOL(val interface{}) *SDK.AttributeValue {
	switch t := val.(type) {
	case bool:
		return &SDK.AttributeValue{BOOL: Boolean(t)}
	}
	return nil
}

func newAttributeValueSS(val interface{}) *SDK.AttributeValue {
	switch t := val.(type) {
	case []string:
		return &SDK.AttributeValue{SS: createPointerSliceString(t)}
	}
	return nil
}

func newAttributeValueNS(val interface{}) *SDK.AttributeValue {
	return &SDK.AttributeValue{NS: MarshalStringSlice(val)}
}

func newAttributeValueBS(val interface{}) *SDK.AttributeValue {
	switch t := val.(type) {
	case [][]byte:
		return &SDK.AttributeValue{BS: t}
	}
	return nil
}

func newAttributeValueM(val interface{}) *SDK.AttributeValue {
	v, ok := val.(map[string]interface{})
	if !ok {
		return nil
	}
	return &SDK.AttributeValue{M: Marshal(v)}
}

func newAttributeValueL(val interface{}) *SDK.AttributeValue {
	// TODO: implement...
	values, ok := val.([]interface{})
	if !ok {
		return nil
	}

	var list []*SDK.AttributeValue
	for _, v := range values {
		list = append(list, createAttributeValue(v))
	}
	return &SDK.AttributeValue{L: list}
}

// Create new AttributeValue from the type of value
func createAttributeValue(v Any) *SDK.AttributeValue {
	switch t := v.(type) {
	case string:
		return &SDK.AttributeValue{
			S: String(t),
		}
	case int, int32, int64, uint, uint32, uint64, float32, float64:
		return &SDK.AttributeValue{
			N: String(fmt.Sprint(t)),
		}
	case []byte:
		return &SDK.AttributeValue{
			B: t,
		}
	case bool:
		return &SDK.AttributeValue{
			BOOL: Boolean(t),
		}
	case []string:
		return &SDK.AttributeValue{
			SS: createPointerSliceString(t),
		}
	case [][]byte:
		return &SDK.AttributeValue{
			BS: t,
		}
	case []int, []int32, []int64, []uint, []uint32, []uint64, []float32, []float64:
		return &SDK.AttributeValue{
			NS: MarshalStringSlice(t),
		}
	case []map[string]interface{}:
		return &SDK.AttributeValue{
			L: createPointerMap(v.([]map[string]interface{})),
		}
	}

	k := reflect.ValueOf(v)
	switch {
	case k.Kind() == reflect.Map:
		return &SDK.AttributeValue{
			M: Marshal(v.(map[string]interface{})),
		}
	}
	return &SDK.AttributeValue{}
}

func createPointerMap(values []map[string]interface{}) []*SDK.AttributeValue {
	var p []*SDK.AttributeValue
	for _, val := range values {
		p = append(p, &SDK.AttributeValue{
			M: Marshal(val),
		})
	}
	return p
}

func createPointerSliceString(values []string) []*string {
	var p []*string
	for _, v := range values {
		str := v
		p = append(p, &str)
	}
	return p
}

// Retrieve value from DynamoDB type
func getItemValue(val *SDK.AttributeValue) Any {
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
		var data []*int
		for _, vString := range val.NS {
			vInt, _ := strconv.Atoi(*vString)
			data = append(data, &vInt)
		}
		return data
	case len(val.SS) > 0:
		var data []*string
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
func Unmarshal(item map[string]*SDK.AttributeValue) map[string]interface{} {
	data := make(map[string]interface{})
	if item == nil {
		return data
	}
	for key, val := range item {
		data[key] = getItemValue(val)
	}
	return data
}

// Convert map to DynamoDb Item data
func Marshal(item map[string]interface{}) map[string]*SDK.AttributeValue {
	data := make(map[string]*SDK.AttributeValue)
	for key, val := range item {
		data[key] = createAttributeValue(val)
	}
	return data
}

// Convert string slice to DynamoDb Item data
func MarshalStringSlice(item Any) []*string {
	var data []*string

	switch reflect.TypeOf(item).Kind() {
	case reflect.Slice:
		val := reflect.ValueOf(item)
		max := val.Len()
		for i := 0; i < max; i++ {
			s := fmt.Sprint(val.Index(i).Interface())
			data = append(data, &s)
		}
	}
	return data
}

func NewProvisionedThroughput(read, write int64) *SDK.ProvisionedThroughput {
	return &SDK.ProvisionedThroughput{
		ReadCapacityUnits:  Long(read),
		WriteCapacityUnits: Long(write),
	}
}

//=======================
//  KeySchema
//=======================

// Create new KeySchema slice
func NewKeySchema(elements ...*SDK.KeySchemaElement) []*SDK.KeySchemaElement {
	if len(elements) > 1 {
		schema := make([]*SDK.KeySchemaElement, 2, 2)
		schema[0] = elements[0]
		schema[1] = elements[1]
		return schema
	} else {
		schema := make([]*SDK.KeySchemaElement, 1, 1)
		schema[0] = elements[0]
		return schema
	}
}

// Create new single KeySchema
func NewKeyElement(keyName, keyType string) *SDK.KeySchemaElement {
	return &SDK.KeySchemaElement{
		AttributeName: String(keyName),
		KeyType:       String(keyType),
	}
}

// Create new single KeySchema for HashKey
func NewHashKeyElement(keyName string) *SDK.KeySchemaElement {
	return NewKeyElement(keyName, KeyTypeHash)
}

// Create new single KeySchema for RangeKey
func NewRangeKeyElement(keyName string) *SDK.KeySchemaElement {
	return NewKeyElement(keyName, KeyTypeRange)
}

//=======================
//  AttributeDefinition
//=======================

// Convert multiple definition to single slice
func NewAttributeDefinitions(attr ...*SDK.AttributeDefinition) []*SDK.AttributeDefinition {
	return attr
}

// Create new definition of table
func NewAttributeDefinition(attrName, attrType string) *SDK.AttributeDefinition {
	newAttr := &SDK.AttributeDefinition{}
	var typ *string
	switch attrType {
	case "S", "N", "B", "BOOL", "L", "M", "SS", "NS", "BS":
		typ = String(attrType)
	default:
		return newAttr
	}
	newAttr.AttributeName = String(attrName)
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
