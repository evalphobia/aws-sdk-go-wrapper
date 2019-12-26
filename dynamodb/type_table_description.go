package dynamodb

import (
	"time"

	SDK "github.com/aws/aws-sdk-go/service/dynamodb"
)

type TableDescription struct {
	// IsExist will be true only when the api response contains table data
	IsExist bool

	ItemCount         int64
	LatestStreamARN   string
	LatestStreamLabel string
	TableARN          string
	TableID           string
	TableName         string
	TableSizeBytes    int64
	TableStatus       string
	CreationDateTime  time.Time

	AttributeDefinitions   []AttributeDefinition
	KeySchema              []KeySchemaElement
	GlobalSecondaryIndexes []GSIDescription
	LocalSecondaryIndexes  []LSIDescription
	ProvisionedThroughput  ProvisionedThroughputDescription
	BillingModeSummary     BillingModeSummary
	RestoreSummary         RestoreSummary
	SSEDescription         SSEDescription
	StreamSpecification    StreamSpecification
}

func NewTableDescription(out *SDK.TableDescription) TableDescription {
	v := TableDescription{}
	if out == nil {
		return v
	}

	v.IsExist = true

	if out.ItemCount != nil {
		v.ItemCount = *out.ItemCount
	}
	if out.LatestStreamArn != nil {
		v.LatestStreamARN = *out.LatestStreamArn
	}
	if out.LatestStreamLabel != nil {
		v.LatestStreamLabel = *out.LatestStreamLabel
	}
	if out.TableArn != nil {
		v.TableARN = *out.TableArn
	}
	if out.TableId != nil {
		v.TableID = *out.TableId
	}
	if out.TableName != nil {
		v.TableName = *out.TableName
	}
	if out.TableSizeBytes != nil {
		v.TableSizeBytes = *out.TableSizeBytes
	}
	if out.TableStatus != nil {
		v.TableStatus = *out.TableStatus
	}
	if out.CreationDateTime != nil {
		v.CreationDateTime = *out.CreationDateTime
	}

	v.AttributeDefinitions = NewAttributeDefinitionList(out.AttributeDefinitions)
	v.BillingModeSummary = NewBillingModeSummary(out.BillingModeSummary)
	v.GlobalSecondaryIndexes = NewGSIDescriptionList(out.GlobalSecondaryIndexes)
	v.KeySchema = NewKeySchemaElementList(out.KeySchema)
	v.LocalSecondaryIndexes = NewLSIDescriptionList(out.LocalSecondaryIndexes)
	v.ProvisionedThroughput = NewProvisionedThroughputDescription(out.ProvisionedThroughput)
	v.RestoreSummary = NewRestoreSummary(out.RestoreSummary)
	v.StreamSpecification = NewStreamSpecification(out.StreamSpecification)
	return v
}

func (d TableDescription) IsEmpty() bool {
	switch {
	case d.IsExist,
		d.ItemCount != 0,
		d.TableSizeBytes != 0,
		d.LatestStreamARN != "",
		d.LatestStreamLabel != "",
		d.TableARN != "",
		d.TableID != "",
		d.TableName != "",
		d.TableStatus != "",
		!d.CreationDateTime.IsZero(),
		len(d.AttributeDefinitions) != 0,
		len(d.KeySchema) != 0:
		return false
	}
	return true
}
