package dynamodb

import (
	SDK "github.com/aws/aws-sdk-go/service/dynamodb"
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
