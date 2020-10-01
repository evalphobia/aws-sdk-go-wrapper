package dynamodb

import (
	SDK "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// BatchGet executes batch_get_item operation
func (t *Table) BatchGet(req BatchGetRequest) (*BatchGetResponse, error) {
	res, err := t.service.client.BatchGetItem(req.ToInput(t.nameWithPrefix))
	if err != nil {
		return nil, err
	}
	return newBatchGetResponse(res, t.nameWithPrefix), nil
}

type BatchGetRequest struct {
	RequestItems           KeysAndAttributes
	ReturnConsumedCapacity string
}

func (r *BatchGetRequest) ToInput(tableName string) *SDK.BatchGetItemInput {
	i := SDK.BatchGetItemInput{}
	i.RequestItems = make(map[string]*SDK.KeysAndAttributes)
	i.RequestItems[tableName] = r.RequestItems.ToSDK()
	if r.ReturnConsumedCapacity != "" {
		i.ReturnConsumedCapacity = &r.ReturnConsumedCapacity
	}
	return &i
}

type BatchGetResponse struct {
	Responses        []map[string]*SDK.AttributeValue
	UnprocessedKeys  KeysAndAttributes
	ConsumedCapacity []ConsumedCapacity
}

func newBatchGetResponse(output *SDK.BatchGetItemOutput, tableName string) *BatchGetResponse {
	r := &BatchGetResponse{}
	if output == nil {
		return r
	}
	if res, ok := output.Responses[tableName]; ok {
		r.Responses = res
	}
	if keys, ok := output.UnprocessedKeys[tableName]; ok {
		r.UnprocessedKeys = newKeysAndAttributes(keys)
	}
	r.ConsumedCapacity = newConsumedCapacities(output.ConsumedCapacity)
	return r
}

func (r BatchGetResponse) Unmarshal(v interface{}) error {
	return r.UnmarshalWithTagName(v, defaultResultTag)
}

// UnmarshalWithTagName unmarshals given slice pointer sturct and tag name from DynamoDB item result to mapping.
func (r BatchGetResponse) UnmarshalWithTagName(v interface{}, structTag string) error {
	decoder := dynamodbattribute.NewDecoder()
	decoder.TagKey = structTag

	items := make([]*SDK.AttributeValue, len(r.Responses))
	for i, m := range r.Responses {
		items[i] = &SDK.AttributeValue{M: m}
	}
	val := &SDK.AttributeValue{
		L: items,
	}
	if err := decoder.Decode(val, v); err != nil {
		return err
	}
	return nil
}
