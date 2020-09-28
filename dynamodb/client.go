package dynamodb

import (
	"fmt"
	"sync"

	SDK "github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
	"github.com/evalphobia/aws-sdk-go-wrapper/private/errors"
	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

const (
	serviceName = "DynamoDB"
)

// DynamoDB has DynamoDB client and table list.
type DynamoDB struct {
	client *SDK.DynamoDB

	logger log.Logger
	prefix string

	tablesMu    sync.RWMutex
	tables      map[string]*Table
	writeTables map[string]struct{}
}

// New returns initializesvc *DynamoDB.
func New(conf config.Config) (*DynamoDB, error) {
	sess, err := conf.Session()
	if err != nil {
		return nil, err
	}

	cli := SDK.New(sess)
	svc := &DynamoDB{
		client:      cli,
		logger:      log.DefaultLogger,
		prefix:      conf.DefaultPrefix,
		tables:      make(map[string]*Table),
		writeTables: make(map[string]struct{}),
	}
	return svc, nil
}

// SetLogger sets logger.
func (svc *DynamoDB) SetLogger(logger log.Logger) {
	svc.logger = logger
}

// ===================
// Table Operation
// ===================

// CreateTable new DynamoDB table
func (svc *DynamoDB) CreateTable(design *TableDesign) error {
	if design.HashKey == nil {
		err := fmt.Errorf("error on `CreateTable`; cannot find hashkey in TableDesign;")
		svc.Errorf("%s", err.Error())
		return err
	}

	originalName := design.name
	design.name = svc.prefix + design.name
	in := design.CreateTableInput()
	out, err := svc.client.CreateTable(in)
	if err != nil {
		svc.Errorf("error on `CreateTable` operation; table=%s; error=%s;", design.GetName(), err.Error())
		design.name = originalName
		return err
	}

	design = newTableDesignFromDescription(NewTableDescription(out.TableDescription))
	svc.Infof("success on `CreateTable` operation; table=%s; status=%s;", design.GetName(), design.status)
	return nil
}

// ForceDeleteTable deletes DynamoDB table.
func (svc *DynamoDB) ForceDeleteTable(name string) error {
	tableName := svc.prefix + name
	in := &SDK.DeleteTableInput{
		TableName: pointers.String(tableName),
	}
	out, err := svc.client.DeleteTable(in)
	if err != nil {
		svc.Errorf("error on `DeleteTable` operation; table=%s; error=%s;", tableName, err.Error())
		return err
	}

	svc.tablesMu.Lock()
	delete(svc.tables, name)
	svc.tablesMu.Unlock()

	desc := NewTableDescription(out.TableDescription)
	svc.Infof("success on `DeleteTable` operation; table=%s; status=%s;", tableName, desc.TableStatus)
	return nil
}

// GetTable returns *Table.
func (svc *DynamoDB) GetTable(name string) (*Table, error) {
	if tbl := svc.GetCachedTable(name); tbl != nil {
		return tbl, nil
	}
	// get the table from AWS api.
	t, err := NewTable(svc, name)
	if err != nil {
		return nil, err
	}

	tableName := svc.prefix + name
	svc.tablesMu.Lock()
	svc.tables[tableName] = t
	svc.tablesMu.Unlock()
	return t, nil
}

// GetCachedTable returns *Table from cache.
func (svc *DynamoDB) GetCachedTable(name string) *Table {
	tableName := svc.prefix + name

	// get the table from cache
	svc.tablesMu.RLock()
	defer svc.tablesMu.RUnlock()
	return svc.tables[tableName]
}

// ListTables gets the list of DynamoDB table.
func (svc *DynamoDB) ListTables() ([]string, error) {
	res, err := svc.client.ListTables(&SDK.ListTablesInput{})
	if err != nil {
		return nil, err
	}

	list := make([]string, len(res.TableNames))
	for i, name := range res.TableNames {
		list[i] = *name
	}
	return list, nil
}

// DescribeTable executes `DescribeTable` operation and get table info.
func (svc *DynamoDB) DescribeTable(name string) (TableDescription, error) {
	res, err := svc.client.DescribeTable(&SDK.DescribeTableInput{
		TableName: pointers.String(name),
	})
	switch {
	case err != nil:
		svc.Errorf("error on `DescribeTable` operation; table=%s; error=%s;", name, err.Error())
		return TableDescription{}, err
	case res == nil:
		err := fmt.Errorf("response is nil")
		svc.Errorf("error on `DescribeTable` operation; table=%s; error=%s;", name, err.Error())
		return TableDescription{}, err
	}

	return NewTableDescription(res.Table), nil
}

// ========================
// Query&Command Operation
// ========================

// PutAll executes put operation for all tables in write spool list.
func (svc *DynamoDB) PutAll() error {
	errList := errors.NewErrors(serviceName)
	for name := range svc.writeTables {
		err := svc.tables[name].Put()
		if err != nil {
			errList.Add(err)
			svc.Errorf("error on `Put` operation; table=%s; error=%s;", name, err.Error())
		}
		svc.removeWriteTable(name)
	}

	if errList.HasError() {
		return errList
	}
	return nil
}

// BatchPutAll executes put operation for all tables with batch operations in write spool list.
func (svc *DynamoDB) BatchPutAll() error {
	errList := errors.NewErrors(serviceName)
	for name := range svc.writeTables {
		err := svc.tables[name].BatchPut()
		if err != nil {
			errList.Add(err)
			svc.Errorf("error on `BatchPut` operation; table=%s; error=%s;", name, err.Error())
		}
		svc.removeWriteTable(name)
	}

	if errList.HasError() {
		return errList
	}
	return nil
}

// addWriteTable adds the table to write spool list.
func (svc *DynamoDB) addWriteTable(tbl *Table) {
	svc.tablesMu.Lock()
	defer svc.tablesMu.Unlock()

	name := tbl.nameWithPrefix
	svc.writeTables[name] = struct{}{}
	if _, ok := svc.tables[name]; !ok {
		svc.tables[name] = tbl
	}
}

// removeWriteTable removes the table from write spool list.
func (svc *DynamoDB) removeWriteTable(name string) {
	svc.tablesMu.Lock()
	defer svc.tablesMu.Unlock()
	delete(svc.writeTables, name)
}

// DoQuery executes `Query` operation and get mapped-items.
func (svc *DynamoDB) DoQuery(in *SDK.QueryInput) (*QueryResult, error) {
	req, err := svc.client.Query(in)
	if err != nil {
		svc.Errorf("error on `Query` operation; table=%s; error=%s", *in.TableName, err.Error())
		return nil, err
	}

	res := &QueryResult{
		Items:            req.Items,
		LastEvaluatedKey: req.LastEvaluatedKey,
		Count:            *req.Count,
		ScannedCount:     *req.ScannedCount,
	}
	return res, nil
}

// BatchGetAll
func (svc *DynamoDB) BatchGetItem(in *SDK.BatchGetItemInput) (*SDK.BatchGetItemOutput, error) {
	return svc.client.BatchGetItem(in)
}

// ========================
// misc
// ========================

// Infof logging information.
func (svc *DynamoDB) Infof(format string, v ...interface{}) {
	svc.logger.Infof(serviceName, format, v...)
}

// Errorf logging error information.
func (svc *DynamoDB) Errorf(format string, v ...interface{}) {
	svc.logger.Errorf(serviceName, format, v...)
}

func newErrors() *errors.Errors {
	return errors.NewErrors(serviceName)
}
