// DynamoDB Table Index

package dynamodb

import (
	DynamoDB "github.com/awslabs/aws-sdk-go/gen/dynamodb"
)

const (
	indexTypeLSI = "local"
	indexTypeGSI = "global"
)

type DynamoIndex struct {
	Name      string
	IndexType string
	KeySchema []DynamoDB.KeySchemaElement
}

func NewDynamoIndex(name, typ string, schema []DynamoDB.KeySchemaElement) *DynamoIndex {
	return &DynamoIndex{
		Name:      name,
		IndexType: typ,
		KeySchema: schema,
	}
}

// get the name of hash key
func (idx *DynamoIndex) GetHashKeyName() string {
	return *idx.KeySchema[0].AttributeName
}

// get the name of range key if exist
func (idx *DynamoIndex) GetRangeKeyName() string {
	if len(idx.KeySchema) > 1 {
		return *idx.KeySchema[1].AttributeName
	} else {
		return ""
	}
}
