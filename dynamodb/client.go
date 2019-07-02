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
	design = newTableDesignFromDescription(out.TableDescription)
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
	if _, ok := svc.tables[name]; ok {
		delete(svc.tables, name)
	}
	svc.tablesMu.Unlock()

	design := newTableDesignFromDescription(out.TableDescription)
	svc.Infof("success on `DeleteTable` operation; table=%s; status=%s;", tableName, design.status)
	return nil
}

// GetTable returns *Table.
func (svc *DynamoDB) GetTable(name string) (*Table, error) {
	tableName := svc.prefix + name

	// get the table from cache
	svc.tablesMu.RLock()
	t, ok := svc.tables[tableName]
	svc.tablesMu.RUnlock()
	if ok {
		return t, nil
	}

	// get the table from AWS api.
	t, err := NewTable(svc, name)
	if err != nil {
		return nil, err
	}

	svc.tablesMu.Lock()
	svc.tables[tableName] = t
	svc.tablesMu.Unlock()
	return t, nil
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
func (svc *DynamoDB) addWriteTable(name string) {
	svc.tablesMu.Lock()
	defer svc.tablesMu.Unlock()
	svc.writeTables[name] = struct{}{}
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
