package dynamodb

import (
	"testing"
)

func TestNewQueryCondition(t *testing.T) {
	c := NewQueryCondition()
	if c.indexName != "" || len(c.conditions) != 0 {
		t.Errorf("error on NewQueryCondition, actual=%v", c)
	}
}

func TestNewCondition(t *testing.T) {
	c := NewCondition(99, "EQ")
	vl := c.AttributeValueList
	if len(vl) != 1 || *c.ComparisonOperator != "EQ" {
		t.Errorf("error on NewCondition, actual=%v", c)
	}
	if *vl[0].N != "99" {
		t.Errorf("error on NewCondition, actual=%v", c)
	}
}

func TestAddEQ(t *testing.T) {
	c := NewQueryCondition()
	c.AddEQ("foo", 99)
	cond := c.conditions
	if _, ok := cond["foo"]; !ok {
		t.Errorf("error on AddEQ, actual=%v", cond)
	}
	vl := cond["foo"].AttributeValueList
	if len(vl) != 1 || *vl[0].N != "99" {
		t.Errorf("error on AddEQ, actual=%v", vl)
	}
}

func TestAddNE(t *testing.T) {
	c := NewQueryCondition()
	c.AddNE("foo", 99)
	cond := c.conditions
	if _, ok := cond["foo"]; !ok {
		t.Errorf("error on AddNE, actual=%v", cond)
	}
	vl := cond["foo"].AttributeValueList
	if len(vl) != 1 || *vl[0].N != "99" {
		t.Errorf("error on AddNE, actual=%v", vl)
	}
}

func TestAddGT(t *testing.T) {
	c := NewQueryCondition()
	c.AddGT("foo", 99)
	cond := c.conditions
	if _, ok := cond["foo"]; !ok {
		t.Errorf("error on AddGT, actual=%v", cond)
	}
	vl := cond["foo"].AttributeValueList
	if len(vl) != 1 || *vl[0].N != "99" {
		t.Errorf("error on AddGT, actual=%v", vl)
	}
}

func TestAddLT(t *testing.T) {
	c := NewQueryCondition()
	c.AddLT("foo", 99)
	cond := c.conditions
	if _, ok := cond["foo"]; !ok {
		t.Errorf("error on AddLT, actual=%v", cond)
	}
	vl := cond["foo"].AttributeValueList
	if len(vl) != 1 || *vl[0].N != "99" {
		t.Errorf("error on AddLT, actual=%v", vl)
	}
}

func TestAddGE(t *testing.T) {
	c := NewQueryCondition()
	c.AddGE("foo", 99)
	cond := c.conditions
	if _, ok := cond["foo"]; !ok {
		t.Errorf("error on AddGE, actual=%v", cond)
	}
	vl := cond["foo"].AttributeValueList
	if len(vl) != 1 || *vl[0].N != "99" {
		t.Errorf("error on AddGE, actual=%v", vl)
	}
}

func TestAddLE(t *testing.T) {
	c := NewQueryCondition()
	c.AddLE("foo", 99)
	cond := c.conditions
	if _, ok := cond["foo"]; !ok {
		t.Errorf("error on AddLE, actual=%v", cond)
	}
	vl := cond["foo"].AttributeValueList
	if len(vl) != 1 || *vl[0].N != "99" {
		t.Errorf("error on AddLE, actual=%v", vl)
	}
}

func TestAdd(t *testing.T) {
	c := NewQueryCondition()
	c.Add("foo", NewCondition(99, "EQ"))
	cond := c.conditions
	if _, ok := cond["foo"]; !ok {
		t.Errorf("error on Add, actual=%v", cond)
	}
	vl := cond["foo"].AttributeValueList
	if len(vl) != 1 || *vl[0].N != "99" {
		t.Errorf("error on Add, actual=%v", vl)
	}
}

func TestUseIndex(t *testing.T) {
	c := NewQueryCondition()
	c.UseIndex("foo-index")
	if c.indexName != "foo-index" {
		t.Errorf("error on UseIndex, actual=%v", c.indexName)
	}
}
