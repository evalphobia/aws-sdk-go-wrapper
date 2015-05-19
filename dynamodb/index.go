// DynamoDB Table Index

package dynamodb

import (
	SDK "github.com/awslabs/aws-sdk-go/service/dynamodb"
)

const (
	indexTypeLSI = "local"
	indexTypeGSI = "global"

	ProjectionTypeAll = "ALL"
	ProjectionTypeKeysOnly = "KEYS_ONLY"
)

// DynamoIndex is wrapper struct for Index,
// used for storing indexes from GetTable's description
type DynamoIndex struct {
	Name      string
	IndexType string
	KeySchema []*SDK.KeySchemaElement
}

// NewDynamoIndex returns initialized DynamoIndex
func NewDynamoIndex(name, typ string, schema []*SDK.KeySchemaElement) *DynamoIndex {
	return &DynamoIndex{
		Name:      name,
		IndexType: typ,
		KeySchema: schema,
	}
}

// GetHashKeyName gets the name of hash key
func (idx *DynamoIndex) GetHashKeyName() string {
	return *idx.KeySchema[0].AttributeName
}

// GetRangeKeyName gets the name of range key if exist
func (idx *DynamoIndex) GetRangeKeyName() string {
	if len(idx.KeySchema) > 1 {
		return *idx.KeySchema[1].AttributeName
	} else {
		return ""
	}
}

// NewLSI returns initilized LocalSecondaryIndex
func NewLSI(name string, schema []*SDK.KeySchemaElement, projection ...string) *SDK.LocalSecondaryIndex {
	lsi := &SDK.LocalSecondaryIndex{
		IndexName: &name,
		KeySchema: schema,
	}

	var proj string
	if len(projection) > 0 {
		proj = projection[0]
	} else {
		proj = ProjectionTypeAll
	}
	lsi.Projection = &SDK.Projection{ProjectionType: &proj}
	return lsi
}

// NewGSI returns initilized GlobalSecondaryIndex
func NewGSI(name string, schema []*SDK.KeySchemaElement, tp *SDK.ProvisionedThroughput, projection ...string) *SDK.GlobalSecondaryIndex {
	gsi := &SDK.GlobalSecondaryIndex{
		IndexName:             &name,
		KeySchema:             schema,
		ProvisionedThroughput: tp,
	}

	var proj string
	if len(projection) > 0 {
		proj = projection[0]
	} else {
		proj = ProjectionTypeAll
	}
	gsi.Projection = &SDK.Projection{ProjectionType: &proj}
	return gsi
}
