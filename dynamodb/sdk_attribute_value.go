// DynamoDB utility

package dynamodb

import (
	"fmt"
	"reflect"
	"strconv"

	SDK "github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

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
	return &SDK.AttributeValue{S: pointers.String(fmt.Sprint(val))}
}

func newAttributeValueN(val interface{}) *SDK.AttributeValue {
	return &SDK.AttributeValue{N: pointers.String(fmt.Sprint(val))}
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
		return &SDK.AttributeValue{BOOL: pointers.Bool(t)}
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
func createAttributeValue(v interface{}) *SDK.AttributeValue {
	switch t := v.(type) {
	case string:
		return &SDK.AttributeValue{
			S: pointers.String(t),
		}
	case int, int32, int64, uint, uint32, uint64, float32, float64:
		return &SDK.AttributeValue{
			N: pointers.String(fmt.Sprint(t)),
		}
	case []byte:
		return &SDK.AttributeValue{
			B: t,
		}
	case bool:
		return &SDK.AttributeValue{
			BOOL: pointers.Bool(t),
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
func getItemValue(val *SDK.AttributeValue) interface{} {
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
		return UnmarshalAttributeValue(val.M)
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
		for _, v := range val.L {
			data = append(data, getItemValue(v))
		}
		return data
	}
	return nil
}

// UnmarshalAttributeValue converts DynamoDB Item to map data.
func UnmarshalAttributeValue(item map[string]*SDK.AttributeValue) map[string]interface{} {
	data := make(map[string]interface{})
	if item == nil {
		return data
	}
	for key, val := range item {
		data[key] = getItemValue(val)
	}
	return data
}

// Marshal converts map data to DynamoDB Item data.
func Marshal(item map[string]interface{}) map[string]*SDK.AttributeValue {
	data := make(map[string]*SDK.AttributeValue)
	for key, val := range item {
		data[key] = createAttributeValue(val)
	}
	return data
}

// MarshalStringSlice converts string slice to DynamoDB Item data.
func MarshalStringSlice(item interface{}) []*string {
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
