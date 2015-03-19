// DynamoDB Client

package dynamodb

import (
	AWS "github.com/awslabs/aws-sdk-go/aws"
	DynamoDB "github.com/awslabs/aws-sdk-go/gen/dynamodb"

	"github.com/evalphobia/aws-sdk-go-wrapper/auth"
	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

const (
	dynamodbConfigSectionName = "dynamodb"
	defaultRegion             = "us-west-1"
	defaultTablePrefix        = "dev_"
)

// wrapped struct for DynamoDB
type AmazonDynamoDB struct {
	client      *DynamoDB.DynamoDB
	tables      map[string]*DynamoTable
	writeTables map[string]bool
}

// Create new AmazonDynamoDB struct
func NewClient() *AmazonDynamoDB {
	d := &AmazonDynamoDB{}
	d.tables = make(map[string]*DynamoTable)
	d.writeTables = make(map[string]bool)
	region := config.GetConfigValue(dynamodbConfigSectionName, "region", defaultRegion)
	d.client = DynamoDB.New(auth.Auth(), region, nil)
	return d
}

// Create new DynamoDB table
func (d *AmazonDynamoDB) CreateTable(in *DynamoDB.CreateTableInput) {
	data, err := d.client.CreateTable(in)
	if err != nil {
		log.Error("[DynamoDB] Error on `CreateTable` operation, table="+*in.TableName, err)
	} else {
		log.Info("[DynamoDB] Complete CreateTable, table="+*in.TableName, data.TableDescription.TableStatus)
	}
}

// get infomation of the table
func (d *AmazonDynamoDB) DescribeTable(name string) (*DynamoDB.TableDescription, error) {
	req, err := d.client.DescribeTable(&DynamoDB.DescribeTableInput{
		TableName: AWS.String(name),
	})
	if err != nil {
		return nil, err
	}
	return req.Table, nil
}

// get the DynamoDB table
func (d *AmazonDynamoDB) GetTable(table string) (*DynamoTable, error) {
	tableName := GetTablePrefix() + table

	// get the table from cache
	t, ok := d.tables[tableName]
	if ok {
		return t, nil
	}

	// get the table info from server
	desc, err := d.DescribeTable(tableName)
	if err != nil {
		return nil, err
	}
	t = &DynamoTable{
		db:      d,
		table:   desc,
		name:    tableName,
		indexes: make(map[string]*DynamoIndex),
	}
	for _, idx := range desc.LocalSecondaryIndexes {
		t.indexes[*idx.IndexName] = NewDynamoIndex(*idx.IndexName, indexTypeLSI, idx.KeySchema)
	}
	for _, idx := range desc.GlobalSecondaryIndexes {
		t.indexes[*idx.IndexName] = NewDynamoIndex(*idx.IndexName, indexTypeGSI, idx.KeySchema)
	}
	d.tables[tableName] = t
	return t, nil
}

// add the table to write spool
func (d *AmazonDynamoDB) addWriteTable(name string) {
	d.writeTables[name] = true
}

// remove the table from write spool
func (d *AmazonDynamoDB) removeWriteTable(name string) {
	d.writeTables[name] = false
}

// execute put operation for all tables in write spool
func (d *AmazonDynamoDB) PutAll() {
	for name, _ := range d.writeTables {
		err := d.tables[name].Put()
		if err != nil {
			log.Error("[DynamoDB] Error on `Put` operation, table="+name, err.Error())
		}
		d.removeWriteTable(name)
	}
}

// get the prefix for DynamoDB table
func GetTablePrefix() string {
	return config.GetConfigValue(dynamodbConfigSectionName, "prefix", defaultTablePrefix)
}

// get the list of DynamoDB table
func (d *AmazonDynamoDB) ListTables() ([]string, error) {
	res, err := d.client.ListTables(&DynamoDB.ListTablesInput{})
	if err != nil {
		return make([]string, 0, 0), err
	}
	return res.TableNames, nil
}
