package dynamodb

import (
	SDK "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// BatchGetItem executes batch_get_item operation
func (svc *DynamoDB) BatchGetItem(in BatchGetItemRequest) (*BatchGetItemResponse, error) {
	o, err := svc.client.BatchGetItem(in.ToInput())
	if err != nil {
		return nil, err
	}
	return newBatchGetItemResponse(o), nil
}

type BatchGetItemRequest struct {
	RequestItems           map[string]KeysAndAttributes
	ReturnConsumedCapacity string
}

func (r *BatchGetItemRequest) ToInput() *SDK.BatchGetItemInput {
	i := SDK.BatchGetItemInput{}
	i.RequestItems = make(map[string]*SDK.KeysAndAttributes, len(r.RequestItems))
	for key, item := range r.RequestItems {
		i.RequestItems[key] = item.ToSDK()
	}
	if r.ReturnConsumedCapacity != "" {
		i.ReturnConsumedCapacity = &r.ReturnConsumedCapacity
	}
	return &i
}

type BatchGetItemResponse struct {
	Responses        map[string][]map[string]*SDK.AttributeValue
	UnprocessedKeys  map[string]KeysAndAttributes
	ConsumedCapacity []ConsumedCapacity `type:"list"`
}

func newBatchGetItemResponse(output *SDK.BatchGetItemOutput) *BatchGetItemResponse {
	r := &BatchGetItemResponse{}
	if output == nil {
		return r
	}
	r.ConsumedCapacity = newConsumedCapacities(output.ConsumedCapacity)
	r.Responses = output.Responses

	r.UnprocessedKeys = make(map[string]KeysAndAttributes, len(output.UnprocessedKeys))
	for key, val := range output.UnprocessedKeys {
		r.UnprocessedKeys[key] = newKeysAndAttributes(val)
	}
	return r
}

// Unmarshal unmarshals given slice pointer sturct from DynamoDB item result to mapping.
//     e.g. err = Unmarshal(&[]*yourStruct)
// The struct tag `dynamodb:""` is used to unmarshal.
func (r BatchGetItemResponse) Unmarshal(v interface{}) error {
	return r.UnmarshalWithTagName(v, defaultResultTag)
}

// UnmarshalWithTagName unmarshals given slice pointer sturct and tag name from DynamoDB item result to mapping.
func (r BatchGetItemResponse) UnmarshalWithTagName(v interface{}, structTag string) error {
	decoder := dynamodbattribute.NewDecoder()
	decoder.TagKey = structTag

	itemsMap := make(map[string]*SDK.AttributeValue, len(r.Responses))
	for tableName, i := range r.Responses {
		items := make([]*SDK.AttributeValue, len(i))
		for i, m := range i {
			items[i] = &SDK.AttributeValue{M: m}
		}
		itemsMap[tableName] = &SDK.AttributeValue{
			L: items,
		}
	}
	val := &SDK.AttributeValue{M: itemsMap}
	if err := decoder.Decode(val, v); err != nil {
		return err
	}
	return nil
}
