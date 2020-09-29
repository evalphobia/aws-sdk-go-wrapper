package dynamodb

import (
	"strconv"
	"time"

	SDK "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

// BillingModeSummary contains the details for the read/write capacity mode.
type BillingModeSummary struct {
	BillingMode               string
	LastUpdateToPayPerRequest time.Time
}

// NewBillingModeSummary creates BillingModeSummary from SDK's output.
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

// IsEmpty checks if the data is empty or not.
func (s BillingModeSummary) IsEmpty() bool {
	switch {
	case s.BillingMode != "",
		!s.LastUpdateToPayPerRequest.IsZero():
		return false
	}
	return true
}

// GSIDescription contains the properties of a global secondary index.
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

// NewGSIDescription creates GSIDescription from SDK's output.
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

// IsEmpty checks if the data is empty or not.
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

// ToGSI converts to SDK's type.
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

// NewGSIDescriptionList creates the list of GSIDescription from SDK's output.
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

// KeySchemaElement represents a single element of a key schema.
type KeySchemaElement struct {
	AttributeName string
	KeyType       string
}

// NewKeySchemaElement creates KeySchemaElement from SDK's output.
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

// IsEmpty checks if the data is empty or not.
func (k KeySchemaElement) IsEmpty() bool {
	switch {
	case k.AttributeName != "",
		k.KeyType != "":
		return false
	}
	return true
}

// ToSDKType converts to SDK's type.
func (k KeySchemaElement) ToSDKType() *SDK.KeySchemaElement {
	if k.IsEmpty() {
		return nil
	}
	return &SDK.KeySchemaElement{
		AttributeName: pointers.String(k.AttributeName),
		KeyType:       pointers.String(k.KeyType),
	}
}

// NewKeySchemaElementList creates the list of KeySchemaElement from SDK's output.
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

// LSIDescription represents the properties of a local secondary index.
type LSIDescription struct {
	IndexARN       string
	IndexName      string
	IndexSizeBytes int64
	ItemCount      int64
	KeySchema      []KeySchemaElement
	Projection     Projection
}

// NewLSIDescription creates LSIDescription from SDK's output.
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

// IsEmpty checks if the data is empty or not.
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

// ToLSI converts to SDK's type.
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

// NewLSIDescriptionList creates the list of LSIDescription from SDK's output.
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

// Projection represents attributes that are copied (projected) from the table into an index.
type Projection struct {
	NonKeyAttributes []string
	ProjectionType   string
}

// NewProjection creates Projection from SDK's output.
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

// IsEmpty checks if the data is empty or not.
func (p Projection) IsEmpty() bool {
	switch {
	case p.ProjectionType != "",
		len(p.NonKeyAttributes) != 0:
		return false
	}
	return true
}

// ToSDKType converts to SDK's type.
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

// ProvisionedThroughputDescription represents the provisioned throughput settings for the table.
type ProvisionedThroughputDescription struct {
	LastDecreaseDateTime   time.Time
	LastIncreaseDateTime   time.Time
	NumberOfDecreasesToday int64
	ReadCapacityUnits      int64
	WriteCapacityUnits     int64
}

// NewProvisionedThroughputDescription creates ProvisionedThroughputDescription from SDK's output.
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

// IsEmpty checks if the data is empty or not.
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

// ToProvisionedThroughput converts to SDK's type.
func (p ProvisionedThroughputDescription) ToProvisionedThroughput() *SDK.ProvisionedThroughput {
	if p.IsEmpty() {
		return nil
	}
	return &SDK.ProvisionedThroughput{
		ReadCapacityUnits:  pointers.Long64(p.ReadCapacityUnits),
		WriteCapacityUnits: pointers.Long64(p.WriteCapacityUnits),
	}
}

// RestoreSummary contains details for the restore.
type RestoreSummary struct {
	RestoreDateTime   time.Time
	RestoreInProgress bool
	SourceBackupARN   string
	SourceTableARN    string
}

// NewRestoreSummary creates RestoreSummary from SDK's output.
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

// IsEmpty checks if the data is empty or not.
func (s RestoreSummary) IsEmpty() bool {
	switch {
	case s.SourceBackupARN != "",
		s.SourceTableARN != "":
		return false
	}
	return true
}

// SSEDescription contains description of the server-side encryption status on the specified table.
type SSEDescription struct {
	KMSMasterKeyARN string
	SSEType         string
	Status          string
}

// NewSSEDescription creates SSEDescription from SDK's output.
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

// IsEmpty checks if the data is empty or not.
func (s SSEDescription) IsEmpty() bool {
	switch {
	case s.KMSMasterKeyARN != "",
		s.SSEType != "",
		s.Status != "":
		return false
	}
	return true
}

// StreamSpecification represents the DynamoDB Streams configuration for a table in DynamoDB.
type StreamSpecification struct {
	StreamEnabled  bool
	StreamViewType string
}

// NewStreamSpecification creates StreamSpecification from SDK's output.
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

// IsEmpty checks if the data is empty or not.
func (s StreamSpecification) IsEmpty() bool {
	switch {
	case s.StreamViewType != "":
		return false
	}
	return true
}

type KeysAndAttributes struct {
	AttributesToGet          []string
	ConsistentRead           bool
	ExpressionAttributeNames map[string]string           `type:"map"`
	Keys                     []map[string]AttributeValue `min:"1" type:"list" required:"true"`
	ProjectionExpression     string                      `type:"string"`
}

func (r KeysAndAttributes) ToSDK() *SDK.KeysAndAttributes {
	o := SDK.KeysAndAttributes{}

	if len(r.Keys) != 0 {
		list := make([]map[string]*SDK.AttributeValue, len(r.Keys))
		for i, v := range r.Keys {
			m := make(map[string]*SDK.AttributeValue, len(v))
			for key, val := range v {
				m[key] = val.ToSDK()
			}
			list[i] = m
		}
		o.Keys = list
	}

	if r.ConsistentRead {
		o.ConsistentRead = pointers.Bool(r.ConsistentRead)
	}
	if r.ProjectionExpression != "" {
		o.ProjectionExpression = pointers.String(r.ProjectionExpression)
	}

	if len(r.AttributesToGet) != 0 {
		attrs := make([]*string, len(r.AttributesToGet))
		for i, s := range r.AttributesToGet {
			attrs[i] = pointers.String(s)
		}
		o.AttributesToGet = attrs
	}
	if len(r.ExpressionAttributeNames) != 0 {
		names := make(map[string]*string, len(r.ExpressionAttributeNames))
		for i, s := range r.ExpressionAttributeNames {
			names[i] = pointers.String(s)
		}
		o.ExpressionAttributeNames = names
	}
	return &o
}

func newKeysAndAttributes(o *SDK.KeysAndAttributes) KeysAndAttributes {
	result := KeysAndAttributes{}
	if o == nil {
		return result
	}

	if len(o.AttributesToGet) != 0 {
		list := make([]string, len(o.AttributesToGet))
		for i, s := range o.AttributesToGet {
			list[i] = *s
		}
		result.AttributesToGet = list
	}

	if len(o.ExpressionAttributeNames) != 0 {
		list := make(map[string]string, len(o.ExpressionAttributeNames))
		for i, s := range o.ExpressionAttributeNames {
			list[i] = *s
		}
		result.ExpressionAttributeNames = list
	}

	if o.ConsistentRead != nil {
		result.ConsistentRead = *o.ConsistentRead
	}
	if o.ProjectionExpression != nil {
		result.ProjectionExpression = *o.ProjectionExpression
	}

	if len(o.Keys) != 0 {
		list := make([]map[string]AttributeValue, len(o.Keys))
		for i, v := range o.Keys {
			list[i] = newAttributeValueMap(v)
		}
		result.Keys = list
	}

	return result
}

type AttributeValue struct {
	Binary         []byte
	BinarySet      [][]byte
	List           []AttributeValue
	Map            map[string]AttributeValue
	Number         string
	NumberInt      int64
	NumberFloat    float64
	NumberSet      []string
	NumberSetInt   []int64
	NumberSetFloat []float64
	Null           bool
	String         string
	StringSet      []string

	Bool      bool
	HasBool   bool
	HasNumber bool
}

func newAttributeValueBySDK(o *SDK.AttributeValue) AttributeValue {
	result := AttributeValue{}
	if o == nil {
		return result
	}

	switch {
	case len(o.B) != 0:
		result.Binary = o.B
	case len(o.BS) != 0:
		result.BinarySet = o.BS
	case len(o.L) != 0:
		result.List = newAttributeValueList(o.L)
	case len(o.M) != 0:
		result.Map = newAttributeValueMap(o.M)
	case o.N != nil:
		result.Number = *o.N
		result.HasNumber = true
	case len(o.NS) != 0:
		list := make([]string, len(o.NS))
		for i, n := range o.NS {
			list[i] = *n
		}
		result.NumberSet = list
	case o.NULL != nil:
		result.Null = *o.NULL
	case o.S != nil:
		result.String = *o.S
	case len(o.SS) != 0:
		list := make([]string, len(o.SS))
		for i, n := range o.SS {
			list[i] = *n
		}
		result.StringSet = list
	case o.BOOL != nil:
		result.Bool = *o.BOOL
		result.HasBool = true
	}
	return result
}

func (r AttributeValue) ToSDK() *SDK.AttributeValue {
	o := SDK.AttributeValue{}

	switch {
	case len(r.Binary) != 0:
		o.B = r.Binary
	case len(r.BinarySet) != 0:
		o.BS = r.BinarySet
	case len(r.List) != 0:
		list := make([]*SDK.AttributeValue, len(r.List))
		for i, v := range r.List {
			list[i] = v.ToSDK()
		}
		o.L = list
	case len(r.Map) != 0:
		m := make(map[string]*SDK.AttributeValue, len(r.Map))
		for key, val := range r.Map {
			m[key] = val.ToSDK()
		}
		o.M = m
	case r.Number != "":
		o.N = pointers.String(r.Number)
	case r.NumberInt != 0:
		o.N = pointers.String(strconv.FormatInt(r.NumberInt, 10))
	case r.NumberFloat != 0:
		o.N = pointers.String(strconv.FormatFloat(r.NumberFloat, 'f', -1, 64))
	case r.HasNumber:
		o.N = pointers.String("0")
	case len(r.NumberSet) != 0:
		list := make([]*string, len(r.NumberSet))
		for i, s := range r.NumberSet {
			list[i] = pointers.String(s)
		}
		o.NS = list
	case len(r.NumberSetInt) != 0:
		list := make([]*string, len(r.NumberSetInt))
		for i, v := range r.NumberSetInt {
			list[i] = pointers.String(strconv.FormatInt(v, 10))
		}
		o.NS = list
	case len(r.NumberSetFloat) != 0:
		list := make([]*string, len(r.NumberSetFloat))
		for i, v := range r.NumberSetFloat {
			list[i] = pointers.String(strconv.FormatFloat(v, 'f', -1, 64))
		}
		o.NS = list
	case r.String != "":
		o.S = pointers.String(r.String)
	case len(r.StringSet) != 0:
		list := make([]*string, len(r.StringSet))
		for i, v := range r.StringSet {
			list[i] = pointers.String(v)
		}
		o.SS = list
	case r.HasBool,
		r.Bool:
		o.BOOL = pointers.Bool(r.Bool)
	case r.Null:
		o.NULL = pointers.Bool(r.Null)
	}
	return &o
}

func newAttributeValueList(list []*SDK.AttributeValue) []AttributeValue {
	if len(list) == 0 {
		return nil
	}

	results := make([]AttributeValue, len(list))
	for i, v := range list {
		results[i] = newAttributeValueBySDK(v)
	}
	return results
}

func newAttributeValueMap(o map[string]*SDK.AttributeValue) map[string]AttributeValue {
	if len(o) == 0 {
		return nil
	}

	m := make(map[string]AttributeValue, len(o))
	for key, val := range o {
		m[key] = newAttributeValueBySDK(val)
	}
	return m
}

type ConsumedCapacity struct {
	CapacityUnits          float64
	GlobalSecondaryIndexes map[string]Capacity
	LocalSecondaryIndexes  map[string]Capacity
	ReadCapacityUnits      float64
	Table                  Capacity
	TableName              string
	WriteCapacityUnits     float64
}

func newConsumedCapacities(list []*SDK.ConsumedCapacity) []ConsumedCapacity {
	if len(list) == 0 {
		return nil
	}

	result := make([]ConsumedCapacity, len(list))
	for i, v := range list {
		result[i] = newConsumedCapacity(v)
	}
	return result
}

func newConsumedCapacity(o *SDK.ConsumedCapacity) ConsumedCapacity {
	result := ConsumedCapacity{}
	if o == nil {
		return result
	}

	if o.CapacityUnits != nil {
		result.CapacityUnits = *o.CapacityUnits
	}
	if o.ReadCapacityUnits != nil {
		result.ReadCapacityUnits = *o.ReadCapacityUnits
	}
	if o.TableName != nil {
		result.TableName = *o.TableName
	}
	if o.WriteCapacityUnits != nil {
		result.WriteCapacityUnits = *o.WriteCapacityUnits
	}

	result.GlobalSecondaryIndexes = newCapacityMap(o.GlobalSecondaryIndexes)
	result.LocalSecondaryIndexes = newCapacityMap(o.LocalSecondaryIndexes)
	result.Table = newCapacity(o.Table)
	return result
}

type Capacity struct {
	CapacityUnits      float64
	ReadCapacityUnits  float64
	WriteCapacityUnits float64
}

func newCapacity(o *SDK.Capacity) Capacity {
	result := Capacity{}
	if o == nil {
		return result
	}

	if o.CapacityUnits != nil {
		result.CapacityUnits = *o.CapacityUnits
	}
	if o.ReadCapacityUnits != nil {
		result.ReadCapacityUnits = *o.ReadCapacityUnits
	}
	if o.WriteCapacityUnits != nil {
		result.WriteCapacityUnits = *o.WriteCapacityUnits
	}
	return result
}

func newCapacityMap(m map[string]*SDK.Capacity) map[string]Capacity {
	if m == nil {
		return nil
	}

	result := make(map[string]Capacity, len(m))
	for key, val := range m {
		result[key] = newCapacity(val)
	}
	return result
}
