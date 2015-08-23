package dynamodb

import (
	"fmt"
	"strings"

	SDK "github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

const (
	conditionEQ      = "="
	conditionLE      = "<="
	conditionLT      = "<"
	conditionGE      = ">="
	conditionGT      = ">"
	conditionBETWEEN = "BETWEEN"

	conditionOR  = "OR"
	conditionAND = "AND"
)

type queryCondition struct {
	Condition string
	Key       string
	Value     interface{}
	SubValue  interface{}
	OR        bool
}

func newQueryCondition(key string, val interface{}) *queryCondition {
	return &queryCondition{
		Key:   key,
		Value: val,
	}
}

func (c *queryCondition) Operator() string {
	if c.OR {
		return conditionOR
	}
	return conditionAND
}

type Query struct {
	table        *DynamoTable
	index        string
	limit        int64
	conditions   map[string]*queryCondition
	isConsistent bool
}

func (q *Query) HasCondition() bool {
	return len(q.conditions) != 0
}

func (q *Query) HasIndex() bool {
	return q.index != ""
}

func (q *Query) HasLimit() bool {
	return q.limit != 0
}

func (q *Query) Limit(i int64) {
	q.limit = i
}

func (q *Query) Index(v string) {
	q.index = v
}

func (q *Query) AndEQ(key string, val interface{}) {
	if _, ok := q.conditions[key]; ok {
		log.Error("[DynamoDB] only contain one condition per key, key="+key, val)
	}
	cond := newQueryCondition(key, val)
	cond.Condition = conditionEQ
	q.conditions[key] = cond
}

func (q *Query) AndLE(key string, val interface{}) {
	if _, ok := q.conditions[key]; ok {
		log.Error("[DynamoDB] only contain one condition per key, key="+key, val)
	}
	cond := newQueryCondition(key, val)
	cond.Condition = conditionLE
	q.conditions[key] = cond
}

func (q *Query) AndLT(key string, val interface{}) {
	if _, ok := q.conditions[key]; ok {
		log.Error("[DynamoDB] only contain one condition per key, key="+key, val)
	}
	cond := newQueryCondition(key, val)
	cond.Condition = conditionLT
	q.conditions[key] = cond
}

func (q *Query) AndGE(key string, val interface{}) {
	if _, ok := q.conditions[key]; ok {
		log.Error("[DynamoDB] only contain one condition per key, key="+key, val)
	}
	cond := newQueryCondition(key, val)
	cond.Condition = conditionGE
	q.conditions[key] = cond
}

func (q *Query) AndGT(key string, val interface{}) {
	if _, ok := q.conditions[key]; ok {
		log.Error("[DynamoDB] only contain one condition per key, key="+key, val)
	}
	cond := newQueryCondition(key, val)
	cond.Condition = conditionGT
	q.conditions[key] = cond
}

func (q *Query) AndBETWEEN(key string, from, to interface{}) {
	if _, ok := q.conditions[key]; ok {
		log.Error("[DynamoDB] only contain one condition per key, key="+key, from)
	}
	cond := newQueryCondition(key, from)
	cond.SubValue = to
	cond.Condition = conditionBETWEEN
	q.conditions[key] = cond
}

func (q *Query) formatConditions() *string {
	var expression []string
	max := len(q.conditions)
	i := 1
	for k, v := range q.conditions {
		var exp string
		switch {
		case v.Condition == conditionBETWEEN:
			exp = fmt.Sprintf("#%s BETWEEN :v_%s AND :s_%s", k, k, k)
		default:
			exp = fmt.Sprintf("#%s %s :v_%s", k, v.Condition, k)
		}
		if i < max {
			exp = exp + " " + v.Operator()
		}
		expression = append(expression, exp)
		i++
	}
	e := strings.Join(expression, " ")
	return &e
}

func (q *Query) formatConditionValues() map[string]*SDK.AttributeValue {
	attrs := q.table.keyAttributes
	m := make(map[string]*SDK.AttributeValue)
	for k, v := range q.conditions {
		typ, ok := attrs[k]
		if !ok {
			continue
		}
		key := fmt.Sprintf(":v_%s", k)
		m[key] = newAttributeValue(typ, v.Value)
		// BETWEEN
		if v.SubValue != nil {
			sub := fmt.Sprintf(":s_%s", k)
			m[sub] = newAttributeValue(typ, v.SubValue)
		}
	}
	return m
}

func (q *Query) formatConditionNames() map[string]*string {
	attrs := q.table.keyAttributes
	m := make(map[string]*string)
	for _, v := range q.conditions {
		if _, ok := attrs[v.Key]; !ok {
			continue
		}
		m[fmt.Sprintf("#%s", v.Key)] = String(v.Key)
	}
	return m
}

func (q *Query) Query() (*QueryResult, error) {
	in := &SDK.QueryInput{
		TableName: String(q.table.name),
	}
	return q.query(in)
}

func (q *Query) Count() (*QueryResult, error) {
	in := &SDK.QueryInput{
		TableName: String(q.table.name),
		Select:    String("COUNT"),
	}
	return q.query(in)
}

func (q *Query) query(in *SDK.QueryInput) (*QueryResult, error) {
	if !q.HasCondition() {
		errData := &DynamoError{}
		errData.AddMessage("condition is missing, you must specify at least one condition")
		return nil, errData
	}
	in.KeyConditionExpression = q.formatConditions()
	in.ExpressionAttributeValues = q.formatConditionValues()
	in.ExpressionAttributeNames = q.formatConditionNames()

	if q.HasIndex() {
		in.IndexName = String(q.index)
	}
	if q.HasLimit() {
		in.Limit = Long(q.limit)
	}

	req, err := q.table.db.client.Query(in)
	if err != nil {
		log.Error("[DynamoDB] Error in `Query` operation, table="+q.table.name, err)
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

type QueryResult struct {
	Items            []map[string]*SDK.AttributeValue
	LastEvaluatedKey map[string]*SDK.AttributeValue
	Count            int64
	ScannedCount     int64
}

func (r QueryResult) ToSliceMap() []map[string]interface{} {
	var m []map[string]interface{}
	for _, item := range r.Items {
		m = append(m, Unmarshal(item))
	}
	return m
}
