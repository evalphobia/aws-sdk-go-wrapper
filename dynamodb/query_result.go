package dynamodb

import (
	SDK "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const defaultResultTag = "dynamodb"

// QueryResult is struct for result of Query operation.
type QueryResult struct {
	Items            []map[string]*SDK.AttributeValue
	LastEvaluatedKey map[string]*SDK.AttributeValue
	Count            int64
	ScannedCount     int64
	tagName          string
}

// ToSliceMap converts result to slice of map.
func (r QueryResult) ToSliceMap() []map[string]interface{} {
	m := make([]map[string]interface{}, len(r.Items))
	for i, item := range r.Items {
		var v map[string]interface{}
		if err := dynamodbattribute.UnmarshalMap(item, &v); err != nil {
			continue
		}
		m[i] = v
	}
	return m
}

// Unmarshal unmarshals given slice pointer sturct from DynamoDB item result to mapping.
//     e.g. err = Unmarshal(&[]*yourStruct)
// The struct tag `dynamodb:""` is used to unmarshal.
func (r QueryResult) Unmarshal(v interface{}) error {
	return r.UnmarshalWithTagName(v, defaultResultTag)
}

// UnmarshalWithTagName unmarshals given slice pointer sturct and tag name from DynamoDB item result to mapping.
func (r QueryResult) UnmarshalWithTagName(v interface{}, structTag string) error {
	decoder := dynamodbattribute.NewDecoder()
	decoder.TagKey = structTag

	items := make([]*SDK.AttributeValue, len(r.Items))
	for i, m := range r.Items {
		items[i] = &SDK.AttributeValue{M: m}
	}
	err := decoder.Decode(&SDK.AttributeValue{L: items}, v)
	return err
}
