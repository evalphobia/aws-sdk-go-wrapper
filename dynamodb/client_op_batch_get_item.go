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
	RequestItems map[string]*SDK.KeysAndAttributes
	ReturnConsumedCapacity string
}

func (r *BatchGetItemRequest) ToInput() *SDK.BatchGetItemInput {
	i := SDK.BatchGetItemInput{
		RequestItems:           r.RequestItems,
	}
	if r.ReturnConsumedCapacity != "" {
		i.ReturnConsumedCapacity = &r.ReturnConsumedCapacity
	}
	return &i
}

type BatchGetItemResponse struct {
	Items map[string][]map[string]*SDK.AttributeValue
	UnprocessedKeys map[string]*SDK.KeysAndAttributes
	ConsumedCapacity []*SDK.ConsumedCapacity `type:"list"`
}

func newBatchGetItemResponse(o *SDK.BatchGetItemOutput) *BatchGetItemResponse {
	res := BatchGetItemResponse{
		Items:           o.Responses,
		UnprocessedKeys: o.UnprocessedKeys,
		ConsumedCapacity: o.ConsumedCapacity,
	}
	return &res
}