package dynamodb

import (
	"time"

	SDK "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

type BillingModeSummary struct {
	BillingMode               string
	LastUpdateToPayPerRequest time.Time
}

func NewBillingModeSummary(out *SDK.BillingModeSummary) BillingModeSummary {
	v := BillingModeSummary{}
	if out == nil {
		return v
	}

	if out.BillingMode != nil {
		v.BillingMode = *out.BillingMode
	}
	if out.LastUpdateToPayPerRequestDateTime != nil {
		v.LastUpdateToPayPerRequest = *out.LastUpdateToPayPerRequestDateTime
	}
	return v
}

func (s BillingModeSummary) IsEmpty() bool {
	switch {
	case s.BillingMode != "",
		!s.LastUpdateToPayPerRequest.IsZero():
		return false
	}
	return true
}

type GSIDescription struct {
	Backfilling           bool
	IndexARN              string
	IndexName             string
	IndexSizeBytes        int64
	IndexStatus           string
	ItemCount             int64
	KeySchema             []KeySchemaElement
	Projection            Projection
	ProvisionedThroughput ProvisionedThroughputDescription
}

func NewGSIDescription(out *SDK.GlobalSecondaryIndexDescription) GSIDescription {
	v := GSIDescription{}
	if out == nil {
		return v
	}

	if out.Backfilling != nil {
		v.Backfilling = *out.Backfilling
	}
	if out.IndexArn != nil {
		v.IndexARN = *out.IndexArn
	}
	if out.IndexName != nil {
		v.IndexName = *out.IndexName
	}
	if out.IndexSizeBytes != nil {
		v.IndexSizeBytes = *out.IndexSizeBytes
	}
	if out.IndexStatus != nil {
		v.IndexStatus = *out.IndexStatus
	}
	if out.ItemCount != nil {
		v.ItemCount = *out.ItemCount
	}

	v.KeySchema = NewKeySchemaElementList(out.KeySchema)
	v.Projection = NewProjection(out.Projection)
	v.ProvisionedThroughput = NewProvisionedThroughputDescription(out.ProvisionedThroughput)
	return v
}

func (d GSIDescription) IsEmpty() bool {
	switch {
	case d.IndexARN != "",
		d.IndexName != "",
		d.IndexStatus != "",
		d.IndexSizeBytes != 0,
		d.ItemCount != 0:
		return false
	}
	return true
}

func (d GSIDescription) ToGSI() *SDK.GlobalSecondaryIndex {
	if d.IsEmpty() {
		return nil
	}

	gsi := &SDK.GlobalSecondaryIndex{
		IndexName:             pointers.String(d.IndexName),
		ProvisionedThroughput: d.ProvisionedThroughput.ToProvisionedThroughput(),
		Projection:            d.Projection.ToSDKType(),
	}
	if len(d.KeySchema) != 0 {
		data := make([]*SDK.KeySchemaElement, 0, len(d.KeySchema))
		for _, s := range d.KeySchema {
			if !s.IsEmpty() {
				data = append(data, s.ToSDKType())
			}
		}
		gsi.KeySchema = data
	}

	return gsi
}

func NewGSIDescriptionList(list []*SDK.GlobalSecondaryIndexDescription) []GSIDescription {
	if len(list) == 0 {
		return nil
	}

	result := make([]GSIDescription, len(list))
	for i, out := range list {
		result[i] = NewGSIDescription(out)
	}
	return result
}

type KeySchemaElement struct {
	AttributeName string
	KeyType       string
}

func NewKeySchemaElement(out *SDK.KeySchemaElement) KeySchemaElement {
	v := KeySchemaElement{}
	if out == nil {
		return v
	}

	if out.AttributeName != nil {
		v.AttributeName = *out.AttributeName
	}
	if out.KeyType != nil {
		v.KeyType = *out.KeyType
	}
	return v
}

func (k KeySchemaElement) IsEmpty() bool {
	switch {
	case k.AttributeName != "",
		k.KeyType != "":
		return false
	}
	return true
}

func (k KeySchemaElement) ToSDKType() *SDK.KeySchemaElement {
	if k.IsEmpty() {
		return nil
	}
	return &SDK.KeySchemaElement{
		AttributeName: pointers.String(k.AttributeName),
		KeyType:       pointers.String(k.KeyType),
	}
}

func NewKeySchemaElementList(list []*SDK.KeySchemaElement) []KeySchemaElement {
	if len(list) == 0 {
		return nil
	}

	result := make([]KeySchemaElement, len(list))
	for i, out := range list {
		result[i] = NewKeySchemaElement(out)
	}
	return result
}

type LSIDescription struct {
	IndexARN       string
	IndexName      string
	IndexSizeBytes int64
	ItemCount      int64
	KeySchema      []KeySchemaElement
	Projection     Projection
}

func NewLSIDescription(out *SDK.LocalSecondaryIndexDescription) LSIDescription {
	v := LSIDescription{}
	if out == nil {
		return v
	}

	if out.IndexArn != nil {
		v.IndexARN = *out.IndexArn
	}
	if out.IndexName != nil {
		v.IndexName = *out.IndexName
	}
	if out.IndexSizeBytes != nil {
		v.IndexSizeBytes = *out.IndexSizeBytes
	}
	if out.ItemCount != nil {
		v.ItemCount = *out.ItemCount
	}

	v.KeySchema = NewKeySchemaElementList(out.KeySchema)
	v.Projection = NewProjection(out.Projection)
	return v
}

func (d LSIDescription) IsEmpty() bool {
	switch {
	case d.IndexARN != "",
		d.IndexName != "",
		d.IndexSizeBytes != 0,
		d.ItemCount != 0:
		return false
	}
	return true
}

func (d LSIDescription) ToLSI() *SDK.LocalSecondaryIndex {
	if d.IsEmpty() {
		return nil
	}

	lsi := &SDK.LocalSecondaryIndex{
		IndexName:  pointers.String(d.IndexName),
		Projection: d.Projection.ToSDKType(),
	}
	if len(d.KeySchema) != 0 {
		data := make([]*SDK.KeySchemaElement, 0, len(d.KeySchema))
		for _, s := range d.KeySchema {
			if !s.IsEmpty() {
				data = append(data, s.ToSDKType())
			}
		}
		lsi.KeySchema = data
	}

	return lsi
}

func NewLSIDescriptionList(list []*SDK.LocalSecondaryIndexDescription) []LSIDescription {
	if len(list) == 0 {
		return nil
	}

	result := make([]LSIDescription, len(list))
	for i, out := range list {
		result[i] = NewLSIDescription(out)
	}
	return result
}

type Projection struct {
	NonKeyAttributes []string
	ProjectionType   string
}

func NewProjection(out *SDK.Projection) Projection {
	v := Projection{}
	if out == nil {
		return v
	}

	if out.ProjectionType != nil {
		v.ProjectionType = *out.ProjectionType
	}

	if len(out.NonKeyAttributes) != 0 {
		data := make([]string, 0, len(out.NonKeyAttributes))
		for _, val := range out.NonKeyAttributes {
			if val != nil {
				data = append(data, *val)
			}
		}
		v.NonKeyAttributes = data
	}
	return v
}

func (p Projection) IsEmpty() bool {
	switch {
	case p.ProjectionType != "",
		len(p.NonKeyAttributes) != 0:
		return false
	}
	return true
}

func (p Projection) ToSDKType() *SDK.Projection {
	if p.IsEmpty() {
		return nil
	}
	pp := &SDK.Projection{
		ProjectionType: pointers.String(p.ProjectionType),
	}
	if len(p.NonKeyAttributes) != 0 {
		data := make([]*string, 0, len(p.NonKeyAttributes))
		for _, v := range p.NonKeyAttributes {
			data = append(data, pointers.String(v))
		}
		pp.NonKeyAttributes = data
	}
	return pp
}

type ProvisionedThroughputDescription struct {
	LastDecreaseDateTime   time.Time
	LastIncreaseDateTime   time.Time
	NumberOfDecreasesToday int64
	ReadCapacityUnits      int64
	WriteCapacityUnits     int64
}

func NewProvisionedThroughputDescription(out *SDK.ProvisionedThroughputDescription) ProvisionedThroughputDescription {
	v := ProvisionedThroughputDescription{}
	if out == nil {
		return v
	}

	if out.LastDecreaseDateTime != nil {
		v.LastDecreaseDateTime = *out.LastDecreaseDateTime
	}
	if out.LastIncreaseDateTime != nil {
		v.LastIncreaseDateTime = *out.LastIncreaseDateTime
	}
	if out.NumberOfDecreasesToday != nil {
		v.NumberOfDecreasesToday = *out.NumberOfDecreasesToday
	}
	if out.ReadCapacityUnits != nil {
		v.ReadCapacityUnits = *out.ReadCapacityUnits
	}
	if out.WriteCapacityUnits != nil {
		v.WriteCapacityUnits = *out.WriteCapacityUnits
	}
	return v
}

func (p ProvisionedThroughputDescription) IsEmpty() bool {
	switch {
	case p.NumberOfDecreasesToday != 0,
		p.ReadCapacityUnits != 0,
		p.WriteCapacityUnits != 0,
		!p.LastDecreaseDateTime.IsZero(),
		!p.LastIncreaseDateTime.IsZero():
		return false
	}
	return true
}

func (p ProvisionedThroughputDescription) ToProvisionedThroughput() *SDK.ProvisionedThroughput {
	if p.IsEmpty() {
		return nil
	}
	return &SDK.ProvisionedThroughput{
		ReadCapacityUnits:  pointers.Long64(p.ReadCapacityUnits),
		WriteCapacityUnits: pointers.Long64(p.WriteCapacityUnits),
	}
}

type RestoreSummary struct {
	RestoreDateTime   time.Time
	RestoreInProgress bool
	SourceBackupARN   string
	SourceTableARN    string
}

func NewRestoreSummary(out *SDK.RestoreSummary) RestoreSummary {
	v := RestoreSummary{}
	if out == nil {
		return v
	}

	if out.RestoreDateTime != nil {
		v.RestoreDateTime = *out.RestoreDateTime
	}
	if out.RestoreInProgress != nil {
		v.RestoreInProgress = *out.RestoreInProgress
	}
	if out.SourceBackupArn != nil {
		v.SourceBackupARN = *out.SourceBackupArn
	}
	if out.SourceTableArn != nil {
		v.SourceTableARN = *out.SourceTableArn
	}
	return v
}

func (s RestoreSummary) IsEmpty() bool {
	switch {
	case s.SourceBackupARN != "",
		s.SourceTableARN != "":
		return false
	}
	return true
}

type SSEDescription struct {
	KMSMasterKeyARN string
	SSEType         string
	Status          string
}

func NewSSEDescription(out *SDK.SSEDescription) SSEDescription {
	v := SSEDescription{}
	if out == nil {
		return v
	}

	if out.KMSMasterKeyArn != nil {
		v.KMSMasterKeyARN = *out.KMSMasterKeyArn
	}
	if out.SSEType != nil {
		v.SSEType = *out.SSEType
	}
	if out.Status != nil {
		v.Status = *out.Status
	}
	return v
}

func (s SSEDescription) IsEmpty() bool {
	switch {
	case s.KMSMasterKeyARN != "",
		s.SSEType != "",
		s.Status != "":
		return false
	}
	return true
}

type StreamSpecification struct {
	StreamEnabled  bool
	StreamViewType string
}

func NewStreamSpecification(out *SDK.StreamSpecification) StreamSpecification {
	v := StreamSpecification{}
	if out == nil {
		return v
	}

	if out.StreamEnabled != nil {
		v.StreamEnabled = *out.StreamEnabled
	}
	if out.StreamViewType != nil {
		v.StreamViewType = *out.StreamViewType
	}
	return v
}

func (s StreamSpecification) IsEmpty() bool {
	switch {
	case s.StreamViewType != "":
		return false
	}
	return true
}
