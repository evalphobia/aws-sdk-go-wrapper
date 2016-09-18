package dynamodb

import (
	SDK "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/mitchellh/mapstructure"
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
		m[i] = UnmarshalAttributeValue(item)
	}
	return m
}

// Unmarshal parse DynamoDB item data and mapping value to given slice pointer sturct.
//     e.g. err = Unmarshal(&[]*yourStruct)
func (r QueryResult) Unmarshal(v interface{}) error {
	m := r.ToSliceMap()

	tagName := r.tagName
	if tagName == "" {
		tagName = defaultResultTag
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   v,
		TagName:  tagName,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(m)
}
