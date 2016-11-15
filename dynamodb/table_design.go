package dynamodb

import (
	"fmt"

	SDK "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

const (
	tableStatusActive = "ACTIVE"
)

// TableDesign is struct for table schema.
type TableDesign struct {
	// for create table
	name          string
	readCapacity  int64
	writeCapacity int64
	HashKey       *SDK.KeySchemaElement
	RangeKey      *SDK.KeySchemaElement
	LSI           []*SDK.LocalSecondaryIndex
	GSI           []*SDK.GlobalSecondaryIndex
	Attributes    map[string]*SDK.AttributeDefinition

	// for table description
	itemCount              int64
	status                 string
	numberOfDecreasesToday int64
}

// ---------------------------------
// construction
// ---------------------------------

// NewTableDesignWithHashKeyS returns create table request data for string hashkey
func NewTableDesignWithHashKeyS(tableName, keyName string) *TableDesign {
	d := newTableDesignWithHashKey(tableName, keyName)
	d.Attributes[keyName] = NewStringAttribute(keyName)
	return d
}

// NewTableDesignWithHashKeyN returns create table request data for number hashkey
func NewTableDesignWithHashKeyN(tableName, keyName string) *TableDesign {
	d := newTableDesignWithHashKey(tableName, keyName)
	d.Attributes[keyName] = NewNumberAttribute(keyName)
	return d
}

func newTableDesignWithHashKey(tableName, hashkeyName string) *TableDesign {
	return &TableDesign{
		name:          tableName,
		HashKey:       NewHashKeyElement(hashkeyName),
		Attributes:    make(map[string]*SDK.AttributeDefinition),
		readCapacity:  1,
		writeCapacity: 1,
	}
}

func newTableDesignFromDescription(desc *SDK.TableDescription) *TableDesign {
	if desc == nil {
		return nil
	}

	d := &TableDesign{
		name:                   *desc.TableName,
		status:                 *desc.TableStatus,
		itemCount:              *desc.ItemCount,
		readCapacity:           *desc.ProvisionedThroughput.ReadCapacityUnits,
		writeCapacity:          *desc.ProvisionedThroughput.WriteCapacityUnits,
		numberOfDecreasesToday: *desc.ProvisionedThroughput.NumberOfDecreasesToday,
		Attributes:             make(map[string]*SDK.AttributeDefinition),
	}
	for _, attr := range desc.AttributeDefinitions {
		d.Attributes[*attr.AttributeName] = attr
	}
	for _, schema := range desc.KeySchema {
		switch *schema.KeyType {
		case "HASH":
			d.HashKey = schema
		case "RANGE":
			d.RangeKey = schema
		}
	}

	for _, lsi := range desc.LocalSecondaryIndexes {
		d.LSI = append(d.LSI, &SDK.LocalSecondaryIndex{
			IndexName:  lsi.IndexName,
			KeySchema:  lsi.KeySchema,
			Projection: lsi.Projection,
		})
	}

	for _, gsi := range desc.GlobalSecondaryIndexes {
		d.GSI = append(d.GSI, &SDK.GlobalSecondaryIndex{
			IndexName:  gsi.IndexName,
			KeySchema:  gsi.KeySchema,
			Projection: gsi.Projection,
			ProvisionedThroughput: &SDK.ProvisionedThroughput{
				ReadCapacityUnits:  gsi.ProvisionedThroughput.ReadCapacityUnits,
				WriteCapacityUnits: gsi.ProvisionedThroughput.WriteCapacityUnits,
			},
		})
	}
	return d
}

// ---------------------------------
// indexes
// ---------------------------------

// AddRangeKeyS adds range key for String type.
func (d *TableDesign) AddRangeKeyS(keyName string) {
	d.RangeKey = NewRangeKeyElement(keyName)
	d.Attributes[keyName] = NewStringAttribute(keyName)
}

// AddRangeKeyN adds range key for Number type.
func (d *TableDesign) AddRangeKeyN(keyName string) {
	d.RangeKey = NewRangeKeyElement(keyName)
	d.Attributes[keyName] = NewNumberAttribute(keyName)
}

// HasRangeKey checks if range key is set or not.
func (d *TableDesign) HasRangeKey() bool {
	return d.RangeKey != nil
}

// HasLSI checks if at least one LocalSecondaryIndex is set or not.
func (d *TableDesign) HasLSI() bool {
	return len(d.LSI) != 0
}

// HasGSI checks if at least one GlobalSecondaryIndex is set or not.
func (d *TableDesign) HasGSI() bool {
	return len(d.GSI) != 0
}

// ListLSI returns multiple LocalSecondaryIndex.
func (d *TableDesign) ListLSI() []*SDK.LocalSecondaryIndex {
	return d.LSI
}

// ListGSI returns multiple GlobalSecondaryIndex.
func (d *TableDesign) ListGSI() []*SDK.GlobalSecondaryIndex {
	return d.GSI
}

// AddLSIS adds LocalSecondaryIndex for String type.
func (d *TableDesign) AddLSIS(name, keyName string) {
	d.Attributes[keyName] = NewStringAttribute(keyName)
	schema := NewKeySchema(d.HashKey, NewRangeKeyElement(keyName))
	lsi := NewLSI(name, schema)
	d.LSI = append(d.LSI, lsi)
}

// AddLSIN adds LocalSecondaryIndex for Number type.
func (d *TableDesign) AddLSIN(name, keyName string) {
	d.Attributes[keyName] = NewNumberAttribute(keyName)
	schema := NewKeySchema(d.HashKey, NewRangeKeyElement(keyName))
	lsi := NewLSI(name, schema)
	d.LSI = append(d.LSI, lsi)
}

func (d *TableDesign) addGSI(name string, key ...string) error {
	var schema []*SDK.KeySchemaElement
	switch len(key) {
	case 1:
		schema = NewKeySchema(NewHashKeyElement(key[0]))
	case 2:
		schema = NewKeySchema(NewHashKeyElement(key[0]), NewRangeKeyElement(key[1]))
	default:
		return fmt.Errorf("keys must have 1 or 2; name=%s; length=%d;", name, len(key))
	}
	tp := newProvisionedThroughput(d.readCapacity, d.writeCapacity)
	gsi := NewGSI(name, schema, tp)
	d.GSI = append(d.GSI, gsi)
	return nil
}

// AddGSIS adds GlobalSecondaryIndex; HashKey=String.
func (d *TableDesign) AddGSIS(name, hashKey string) error {
	d.Attributes[hashKey] = NewStringAttribute(hashKey)
	return d.addGSI(name, hashKey)
}

// AddGSIN adds GlobalSecondaryIndex; HashKey=Number.
func (d *TableDesign) AddGSIN(name, hashKey string) error {
	d.Attributes[hashKey] = NewNumberAttribute(hashKey)
	return d.addGSI(name, hashKey)
}

// AddGSISS adds GlobalSecondaryIndex; HashKey=String, RangeKey=String.
func (d *TableDesign) AddGSISS(name, hashKey, rangeKey string) error {
	d.Attributes[hashKey] = NewStringAttribute(hashKey)
	d.Attributes[rangeKey] = NewStringAttribute(rangeKey)
	return d.addGSI(name, hashKey, rangeKey)
}

// AddGSISN adds GlobalSecondaryIndex; HashKey=String, RangeKey=Number.
func (d *TableDesign) AddGSISN(name, hashKey, rangeKey string) error {
	d.Attributes[hashKey] = NewStringAttribute(hashKey)
	d.Attributes[rangeKey] = NewNumberAttribute(rangeKey)
	return d.addGSI(name, hashKey, rangeKey)
}

// AddGSINN adds GlobalSecondaryIndex; HashKey=Number, RangeKey=Number.
func (d *TableDesign) AddGSINN(name, hashKey, rangeKey string) error {
	d.Attributes[hashKey] = NewNumberAttribute(hashKey)
	d.Attributes[rangeKey] = NewNumberAttribute(rangeKey)
	return d.addGSI(name, hashKey, rangeKey)
}

// AddGSINS adds GlobalSecondaryIndex; HashKey=Number, RangeKey=String.
func (d *TableDesign) AddGSINS(name, hashKey, rangeKey string) error {
	d.Attributes[hashKey] = NewNumberAttribute(hashKey)
	d.Attributes[rangeKey] = NewStringAttribute(rangeKey)
	return d.addGSI(name, hashKey, rangeKey)
}

// ---------------------------------
// Attributes
// ---------------------------------

// AttributeList returns list of *SDK.AttributeDefinition.
func (d *TableDesign) AttributeList() []*SDK.AttributeDefinition {
	var attrs []*SDK.AttributeDefinition
	for _, v := range d.Attributes {
		attrs = append(attrs, v)
	}
	return attrs
}

// GetKeyAttributes returns KeyAttributes.
func (d *TableDesign) GetKeyAttributes() map[string]string {
	m := make(map[string]string)
	for _, attr := range d.Attributes {
		m[*attr.AttributeName] = *attr.AttributeType
	}
	return m
}

// ---------------------------------
// Throughput
// ---------------------------------

// SetThroughput sets read and write throughput.
func (d *TableDesign) SetThroughput(r, w int64) {
	d.readCapacity = r
	d.writeCapacity = w
}

// GetReadCapacity returns read capacity.
func (d *TableDesign) GetReadCapacity() int64 {
	return d.readCapacity
}

// GetWriteCapacity returns write capacity.
func (d *TableDesign) GetWriteCapacity() int64 {
	return d.writeCapacity
}

// GetNumberOfDecreasesToday returns NumberOfDecreasesToday for throughput.
func (d *TableDesign) GetNumberOfDecreasesToday() int64 {
	return d.numberOfDecreasesToday
}

// ---------------------------------
// misc
// ---------------------------------

// CreateTableInput creates *SDK.CreateTableInput from the table design.
func (d *TableDesign) CreateTableInput(prefix string) *SDK.CreateTableInput {
	var keys []*SDK.KeySchemaElement
	keys = append(keys, d.HashKey)
	if d.HasRangeKey() {
		keys = append(keys, d.RangeKey)
	}

	in := &SDK.CreateTableInput{
		TableName:             pointers.String(prefix + d.name),
		KeySchema:             keys,
		AttributeDefinitions:  d.AttributeList(),
		ProvisionedThroughput: newProvisionedThroughput(d.readCapacity, d.writeCapacity),
	}

	if d.HasLSI() {
		in.LocalSecondaryIndexes = d.ListLSI()
	}
	if d.HasGSI() {
		in.GlobalSecondaryIndexes = d.ListGSI()
	}
	return in
}

// keyAttributeValue returns map[string]*SDK.CreateTableInput for Key.
func (d *TableDesign) keyAttributeValue(hashValue interface{}, rangeValue ...interface{}) map[string]*SDK.AttributeValue {
	key := make(map[string]*SDK.AttributeValue)
	key[d.GetHashKeyName()] = createAttributeValue(hashValue)

	rangeKey := d.GetRangeKeyName()
	if len(rangeValue) == 1 && rangeKey != "" {
		key[rangeKey] = createAttributeValue(rangeValue[0])
	}
	return key
}

// GetName returns table name.
func (d *TableDesign) GetName() string {
	return d.name
}

// GetStatus returns table status.
func (d *TableDesign) GetStatus() string {
	return d.status
}

// IsActive checks if the table status is active or not.
func (d *TableDesign) IsActive() bool {
	return d.status == tableStatusActive
}

// GetItemCount returns items count on this table.
func (d *TableDesign) GetItemCount() int64 {
	return d.itemCount
}

// GetHashKeyName returns attribute name of the HashKey.
func (d *TableDesign) GetHashKeyName() string {
	return *d.HashKey.AttributeName
}

// GetRangeKeyName returns attribute name of the RangeKey.
func (d *TableDesign) GetRangeKeyName() string {
	if !d.HasRangeKey() {
		return ""
	}

	return *d.RangeKey.AttributeName
}
