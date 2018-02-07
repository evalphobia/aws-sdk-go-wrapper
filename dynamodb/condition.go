package dynamodb

import (
	"fmt"
	"strings"

	SDK "github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
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

// ConditionList contains multiple condition.
type ConditionList struct {
	keyAttributes map[string]string
	conditions    map[string]*Condition
	filters       map[string]*Condition

	index        string
	limit        int64
	startKey     map[string]*SDK.AttributeValue
	isConsistent bool
	isDesc       bool // descending order
}

// NewConditionList returns initialized *ConditionList.
func NewConditionList(keyAttributes map[string]string) *ConditionList {
	return &ConditionList{
		keyAttributes: keyAttributes,
		conditions:    make(map[string]*Condition),
		filters:       make(map[string]*Condition),
	}
}

// HasCondition checks if at least one condition is set or not.
func (c *ConditionList) HasCondition() bool {
	return len(c.conditions) != 0
}

// HasFilter checks if at least one filter is set or not.
func (c *ConditionList) HasFilter() bool {
	return len(c.filters) != 0
}

// HasIndex checks if the index is set or not.
func (c *ConditionList) HasIndex() bool {
	return c.index != ""
}

// HasLimit checks if limit number is set or not.
func (c *ConditionList) HasLimit() bool {
	return c.limit != 0
}

// SetLimit sets limit number.
func (c *ConditionList) SetLimit(i int64) {
	c.limit = i
}

// SetIndex sets index to use.
func (c *ConditionList) SetIndex(v string) {
	c.index = v
}

// SetConsistent sets consistent read flag.
func (c *ConditionList) SetConsistent(b bool) {
	c.isConsistent = b
}

// SetDesc sets descending order flag.
func (c *ConditionList) SetDesc(b bool) {
	c.isDesc = b
}

// SetStartKey sets ExclusiveStartKey.
func (c *ConditionList) SetStartKey(startKey map[string]*SDK.AttributeValue) {
	c.startKey = startKey
}

// AndEQ adds EQ(equal) condition.
func (c *ConditionList) AndEQ(key string, val interface{}) {
	c.setCondition(conditionEQ, key, val)
}

// AndLE adds LE(less equal than) condition.
func (c *ConditionList) AndLE(key string, val interface{}) {
	c.setCondition(conditionLE, key, val)
}

// AndLT adds LT(less than) condition.
func (c *ConditionList) AndLT(key string, val interface{}) {
	c.setCondition(conditionLT, key, val)
}

// AndGE adds GE(greater equal than) condition.
func (c *ConditionList) AndGE(key string, val interface{}) {
	c.setCondition(conditionGE, key, val)
}

// AndGT adds GT(greater than) condition.
func (c *ConditionList) AndGT(key string, val interface{}) {
	c.setCondition(conditionGT, key, val)
}

// AndBETWEEN adds BETWEEN condition.
func (c *ConditionList) AndBETWEEN(key string, from, to interface{}) {
	c.setCondition(conditionBETWEEN, key, from, to)
}

func (c *ConditionList) setCondition(condition, key string, val interface{}, subVal ...interface{}) {
	if _, ok := c.conditions[key]; ok {
		return
	}

	cond := newCondition(condition, key, val)
	if len(subVal) == 1 {
		cond.SubValue = subVal[0]
	}
	c.conditions[key] = cond
}

// FilterEQ adds EQ(equal) filter.
func (c *ConditionList) FilterEQ(key string, val interface{}) {
	c.setFilter(conditionEQ, key, val)
}

// FilterLE adds LE(less equal than) filter.
func (c *ConditionList) FilterLE(key string, val interface{}) {
	c.setFilter(conditionLE, key, val)
}

// FilterLT adds LT(less than) filter.
func (c *ConditionList) FilterLT(key string, val interface{}) {
	c.setFilter(conditionLT, key, val)
}

// FilterGE adds GE(greater equal than) filter.
func (c *ConditionList) FilterGE(key string, val interface{}) {
	c.setFilter(conditionGE, key, val)
}

// FilterGT adds GT(greater than) filter.
func (c *ConditionList) FilterGT(key string, val interface{}) {
	c.setFilter(conditionGT, key, val)
}

// FilterBETWEEN adds BETWEEN filter.
func (c *ConditionList) FilterBETWEEN(key string, from, to interface{}) {
	c.setFilter(conditionBETWEEN, key, from, to)
}

func (c *ConditionList) setFilter(condition, key string, val interface{}, subVal ...interface{}) {
	if _, ok := c.filters[key]; ok {
		return
	}

	cond := newCondition(condition, key, val)
	cond.isFilter = true
	if len(subVal) == 1 {
		cond.SubValue = subVal[0]
	}
	c.filters[key] = cond
}

// FormatCondition returns string pointer for KeyConditionExpression.
func (c *ConditionList) FormatCondition() *string {
	return c.formatCondition(c.conditions)
}

// FormatFilter returns string pointer for KeyConditionExpression.
func (c *ConditionList) FormatFilter() *string {
	return c.formatCondition(c.filters)
}

// formatCondition returns string pointer for ConditionExpression and FilterExpression.
func (c *ConditionList) formatCondition(conditions map[string]*Condition) *string {
	max := len(conditions)
	if max == 0 {
		return nil
	}

	i := 1
	expression := make([]string, 0, max)
	for key, cond := range conditions {
		exp := cond.expression(key)

		// add space unless final expression
		if i < max {
			exp = exp + " " + cond.operator()
		}
		expression = append(expression, exp)
		i++
	}
	e := strings.Join(expression, " ")
	return &e
}

// FormatValues returns the parameter for ExpressionAttributeValues.
func (c *ConditionList) FormatValues() map[string]*SDK.AttributeValue {
	attrs := c.keyAttributes
	m := make(map[string]*SDK.AttributeValue)

	for k, cond := range c.getMergedConditions() {
		typ, ok := attrs[k]
		if !ok {
			continue
		}
		key := cond.valueName()
		m[key] = newAttributeValue(typ, cond.Value)
		// BETWEEN
		if cond.SubValue != nil {
			sub := cond.subValueName()
			m[sub] = newAttributeValue(typ, cond.SubValue)
		}
	}
	return m
}

// FormatNames returns the parameter for ExpressionAttributeNames.
func (c *ConditionList) FormatNames() map[string]*string {
	attrs := c.keyAttributes
	m := make(map[string]*string)

	for _, cond := range c.getMergedConditions() {
		if _, ok := attrs[cond.Key]; !ok {
			continue
		}
		m[cond.keyName()] = pointers.String(cond.Key)
	}
	return m
}

func (c *ConditionList) getMergedConditions() map[string]*Condition {
	list := make(map[string]*Condition)
	for k, v := range c.filters {
		list[k] = v
	}

	for k, v := range c.conditions {
		list[k] = v
	}
	return list
}

// Condition contains condition.
type Condition struct {
	Condition string
	Key       string
	Value     interface{}
	SubValue  interface{}
	OR        bool
	isFilter  bool
}

// newCondition returns initialized *Condition.
func newCondition(condition, key string, val interface{}) *Condition {
	return &Condition{
		Condition: condition,
		Key:       key,
		Value:     val,
	}
}

func (c *Condition) expression(key string) (expression string) {
	switch {
	case c.Condition == conditionBETWEEN:
		return fmt.Sprintf("%s BETWEEN %s AND %s", c.keyName(), c.valueName(), c.subValueName())
	default:
		return fmt.Sprintf("%s %s %s", c.keyName(), c.Condition, c.valueName())
	}
}

func (c *Condition) operator() string {
	if c.OR {
		return conditionOR
	}
	return conditionAND
}

func (c *Condition) keyName() string {
	switch {
	case c.isFilter:
		return fmt.Sprintf("#f_%s", c.Key)
	default:
		return fmt.Sprintf("#c_%s", c.Key)
	}
}

func (c *Condition) valueName() string {
	switch {
	case c.isFilter:
		return fmt.Sprintf(":f_%s", c.Key)
	default:
		return fmt.Sprintf(":c_%s", c.Key)
	}
}

func (c *Condition) subValueName() string {
	switch {
	case c.isFilter:
		return fmt.Sprintf(":fs_%s", c.Key)
	default:
		return fmt.Sprintf(":cs_%s", c.Key)
	}
}
