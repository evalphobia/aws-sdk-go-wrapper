// DynamoDB Table Index

package dynamodb

import (
	SDK "github.com/aws/aws-sdk-go/service/dynamodb"
)

// projection type of the index.
// see: http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_Projection.html
const (
	ProjectionTypeAll      = "ALL"
	ProjectionTypeKeysOnly = "KEYS_ONLY"
)

// NewLSI returns initilized LocalSecondaryIndex.
func NewLSI(name string, schema []*SDK.KeySchemaElement, projection ...string) *SDK.LocalSecondaryIndex {
	var proj string
	switch {
	case len(projection) == 1:
		proj = projection[0]
	default:
		proj = ProjectionTypeAll
	}

	return &SDK.LocalSecondaryIndex{
		IndexName:  &name,
		KeySchema:  schema,
		Projection: &SDK.Projection{ProjectionType: &proj},
	}
}

// NewGSI returns initilized GlobalSecondaryIndex.
func NewGSI(name string, schema []*SDK.KeySchemaElement, tp *SDK.ProvisionedThroughput, projection ...string) *SDK.GlobalSecondaryIndex {
	var proj string
	switch {
	case len(projection) == 1:
		proj = projection[0]
	default:
		proj = ProjectionTypeAll
	}

	return &SDK.GlobalSecondaryIndex{
		IndexName:             &name,
		KeySchema:             schema,
		ProvisionedThroughput: tp,
		Projection:            &SDK.Projection{ProjectionType: &proj},
	}
}
