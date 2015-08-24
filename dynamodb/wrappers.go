// misculious wrapper structs

package dynamodb

import (
	SDK "github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

const (
	tableStatusActive = "ACTIVE"
)

// TableDescription is wrapper struct fro SDK.TableDescription
type TableDescription struct {
	*SDK.TableDescription
}

func (d TableDescription) GetItemCount() int64 {
	return *d.ItemCount
}

func (d TableDescription) GetTableName() string {
	return *d.TableName
}

func (d TableDescription) GetTableStatus() string {
	return *d.TableStatus
}

func (d TableDescription) IsActive() bool {
	return d.GetTableStatus() == tableStatusActive
}

func (d TableDescription) GetReadCapacityUnits() int64 {
	return *d.ProvisionedThroughput.ReadCapacityUnits
}

func (d TableDescription) GetWriteCapacityUnits() int64 {
	return *d.ProvisionedThroughput.WriteCapacityUnits
}

func (d TableDescription) GetNumberOfDecreasesToday() int64 {
	return *d.ProvisionedThroughput.NumberOfDecreasesToday
}

func (d TableDescription) GetKeyAttributes() map[string]string {
	m := make(map[string]string)
	for _, attr := range d.AttributeDefinitions {
		m[*attr.AttributeName] = *attr.AttributeType
	}
	return m
}

// CreateTableInput is wrapper struct for CreateTable operation
type CreateTableInput struct {
	Name          string
	HashKey       *SDK.KeySchemaElement
	RangeKey      *SDK.KeySchemaElement
	LSI           []*SDK.LocalSecondaryIndex
	GSI           []*SDK.GlobalSecondaryIndex
	ReadCapacity  int64
	WriteCapacity int64
	Attributes    map[string]*SDK.AttributeDefinition
}

func newCreateTableWithHashKey(tableName, hashkeyName string) *CreateTableInput {
	return &CreateTableInput{
		Name:          tableName,
		HashKey:       NewHashKeyElement(hashkeyName),
		Attributes:    make(map[string]*SDK.AttributeDefinition),
		ReadCapacity:  1,
		WriteCapacity: 1,
	}
}

// NewCreateTableWithHashKeyS returns create table request data for string hashkey
func NewCreateTableWithHashKeyS(tableName, keyName string) *CreateTableInput {
	ct := newCreateTableWithHashKey(tableName, keyName)
	ct.Attributes[keyName] = NewStringAttribute(keyName)
	return ct
}

// NewCreateTableWithHashKeyN returns create table request data for number hashkey
func NewCreateTableWithHashKeyN(tableName, keyName string) *CreateTableInput {
	ct := newCreateTableWithHashKey(tableName, keyName)
	ct.Attributes[keyName] = NewNumberAttribute(keyName)
	return ct
}

func (ct *CreateTableInput) AddRangeKeyS(keyName string) {
	ct.RangeKey = NewRangeKeyElement(keyName)
	ct.Attributes[keyName] = NewStringAttribute(keyName)
}

func (ct *CreateTableInput) AddRangeKeyN(keyName string) {
	ct.RangeKey = NewRangeKeyElement(keyName)
	ct.Attributes[keyName] = NewNumberAttribute(keyName)
}

func (ct *CreateTableInput) HasRangeKey() bool {
	return ct.RangeKey != nil
}

func (ct *CreateTableInput) HasLSI() bool {
	return len(ct.LSI) != 0
}

func (ct *CreateTableInput) HasGSI() bool {
	return len(ct.GSI) != 0
}

func (ct *CreateTableInput) ListLSI() []*SDK.LocalSecondaryIndex {
	return ct.LSI
}

func (ct *CreateTableInput) ListGSI() []*SDK.GlobalSecondaryIndex {
	return ct.GSI
}

func (ct *CreateTableInput) AddLSIS(name, keyName string) {
	ct.Attributes[keyName] = NewStringAttribute(keyName)
	schema := NewKeySchema(ct.HashKey, NewRangeKeyElement(keyName))
	lsi := NewLSI(name, schema)
	ct.LSI = append(ct.LSI, lsi)
}

func (ct *CreateTableInput) AddLSIN(name, keyName string) {
	ct.Attributes[keyName] = NewNumberAttribute(keyName)
	schema := NewKeySchema(ct.HashKey, NewRangeKeyElement(keyName))
	lsi := NewLSI(name, schema)
	ct.LSI = append(ct.LSI, lsi)
}

func (ct *CreateTableInput) addGSI(name string, key ...string) {
	var schema []*SDK.KeySchemaElement
	switch len(key) {
	case 1:
		schema = NewKeySchema(NewHashKeyElement(key[0]))
	case 2:
		schema = NewKeySchema(NewHashKeyElement(key[0]), NewRangeKeyElement(key[1]))
	default:
		log.Error("[DynamoDB] keys must have 1 or 2, name="+name, len(key))
		return
	}
	tp := NewProvisionedThroughput(ct.ReadCapacity, ct.WriteCapacity)
	gsi := NewGSI(name, schema, tp)
	ct.GSI = append(ct.GSI, gsi)
}

func (ct *CreateTableInput) AddGSIS(name, hashKey string) {
	ct.Attributes[hashKey] = NewStringAttribute(hashKey)
	ct.addGSI(name, hashKey)
}

func (ct *CreateTableInput) AddGSIN(name, hashKey string) {
	ct.Attributes[hashKey] = NewNumberAttribute(hashKey)
	ct.addGSI(name, hashKey)
}

func (ct *CreateTableInput) AddGSISS(name, hashKey, rangeKey string) {
	ct.Attributes[hashKey] = NewStringAttribute(hashKey)
	ct.Attributes[rangeKey] = NewStringAttribute(rangeKey)
	ct.addGSI(name, hashKey, rangeKey)
}

func (ct *CreateTableInput) AddGSISN(name, hashKey, rangeKey string) {
	ct.Attributes[hashKey] = NewStringAttribute(hashKey)
	ct.Attributes[rangeKey] = NewNumberAttribute(rangeKey)
	ct.addGSI(name, hashKey, rangeKey)
}

func (ct *CreateTableInput) AddGSINN(name, hashKey, rangeKey string) {
	ct.Attributes[hashKey] = NewNumberAttribute(hashKey)
	ct.Attributes[rangeKey] = NewNumberAttribute(rangeKey)
	ct.addGSI(name, hashKey, rangeKey)
}

func (ct *CreateTableInput) AddGSINS(name, hashKey, rangeKey string) {
	ct.Attributes[hashKey] = NewNumberAttribute(hashKey)
	ct.Attributes[rangeKey] = NewStringAttribute(rangeKey)
	ct.addGSI(name, hashKey, rangeKey)
}

func (ct *CreateTableInput) AttributeList() []*SDK.AttributeDefinition {
	var attrs []*SDK.AttributeDefinition
	for _, v := range ct.Attributes {
		attrs = append(attrs, v)
	}
	return attrs
}

func (ct *CreateTableInput) SetThroughput(r, w int64) {
	ct.ReadCapacity = r
	ct.WriteCapacity = w
}
