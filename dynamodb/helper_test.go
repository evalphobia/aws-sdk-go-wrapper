package dynamodb

import (
	"bytes"
	"testing"

	"fmt"
)

var _ = fmt.Sprint("")

func TestCreateAttributeValue(t *testing.T) {
	s := createAttributeValue("foo")
	if *s.S != "foo" {
		t.Errorf("error on createAttributeValue, actual=%+v", s)
	}

	i := createAttributeValue(99)
	if *i.N != "99" {
		t.Errorf("error on createAttributeValue, actual=%+v", i)
	}

	b := createAttributeValue([]byte("foo"))
	if !bytes.Equal(b.B, []byte{102, 111, 111}) {
		t.Errorf("error on createAttributeValue, actual=%+v", b)
	}

	bl := createAttributeValue(true)
	if *bl.BOOL != true {
		t.Errorf("error on createAttributeValue, actual=%+v", bl)
	}

	ssData := []string{"foo1", "foo2", "foo3"}
	ss := createAttributeValue(ssData)
	if len(ss.SS) != 3 || *ss.SS[0] != "foo1" {
		t.Errorf("error on createAttributeValue, actual=%s, values=%+v", *ss.SS[0], ss)
	}

	nsData := []int{1, 2, 3, 9}
	ns := createAttributeValue(nsData)
	if len(ns.NS) != 4 || *ns.NS[0] != "1" {
		t.Errorf("error on createAttributeValue, actual=%+v", ns)
	}

	var bsData [][]byte
	bsData = append(bsData, []byte("bs1"), []byte("bs2"), []byte("bs3"))
	bs := createAttributeValue(bsData)
	if len(bs.BS) != 3 {
		t.Errorf("error on createAttributeValue, actual=%+v", bs)
	}

	mData := make(map[string]interface{})
	mData["id"] = 1
	mData["data"] = "foo"
	m := createAttributeValue(mData)
	if val, ok := (*m.M)["id"]; !ok || getItemValue(val) != 1 {
		t.Errorf("error on createAttributeValue, actual=%+v", m)
	}

	st := createAttributeValue(TestStruct{})
	if st.S != nil || st.N != nil || st.BOOL != nil || len(st.B) != 0 {
		t.Errorf("error on createAttributeValue, actual=%+v", st)
	}
}

func TestGetItemValue(t *testing.T) {
	s := createAttributeValue("foo")
	if getItemValue(s) != "foo" {
		t.Errorf("error on getItemValue, actual=%v", s)
	}

	i := createAttributeValue(99)
	if getItemValue(i) != 99 {
		t.Errorf("error on getItemValue, actual=%v", i)
	}

	b := createAttributeValue([]byte("foo"))
	if !bytes.Equal(getItemValue(b).([]byte), []byte{102, 111, 111}) {
		t.Errorf("error on getItemValue, actual=%v", b)
	}

	bl := createAttributeValue(true)
	if getItemValue(bl) != true {
		t.Errorf("error on getItemValue, actual=%v", bl)
	}

	ns := createAttributeValue([]int{1, 2, 3})
	nsValue := getItemValue(ns).([]*int)
	if len(nsValue) != 3 || *nsValue[0] != 1 || *nsValue[1] != 2 || *nsValue[2] != 3 {
		t.Errorf("error on getItemValue, actual=%+v", nsValue)
	}

	ss := createAttributeValue([]string{"foo1", "foo2", "foo3"})
	ssValue := getItemValue(ss).([]*string)
	if len(ssValue) != 3 || *ssValue[0] != "foo1" || *ssValue[1] != "foo2" || *ssValue[2] != "foo3" {
		t.Errorf("error on getItemValue, actual=%+v", ssValue)
	}

	var bsData [][]byte
	bsData = append(bsData, []byte("bs1"), []byte("bs2"), []byte("bs3"))
	bs := createAttributeValue(bsData)
	bsValue := getItemValue(bs).([][]byte)
	if len(bsValue) != 3 || !bytes.Equal(bsValue[0], []byte("bs1")) ||
		!bytes.Equal(bsValue[1], []byte("bs2")) || !bytes.Equal(bsValue[2], []byte("bs3")) {
		t.Errorf("error on createAttributeValue, actual=%+v", bsValue)
	}

	mData := make(map[string]interface{})
	mData["id"] = 1
	mData["data"] = "foo"
	m := createAttributeValue(mData)
	mValue := getItemValue(m).(map[string]interface{})
	if val, ok := mValue["id"]; !ok || val != 1 {
		t.Errorf("error on createAttributeValue, actual=%+v", mValue)
	}

}

// TestUnmarshal TODO: write test
func TestUnmarshal(t *testing.T) {
	t.Skip("TODO: write test")
}

func TestNewProvisionedThroughput(t *testing.T) {
	tp := NewProvisionedThroughput(80, 600)
	if *tp.ReadCapacityUnits != 80 || *tp.WriteCapacityUnits != 600 {
		t.Errorf("error on NewProvisionedThroughput, actual=%v", tp)
	}
}

func TestNewKeySchema(t *testing.T) {
	key := NewKeyElement("foo", "bar")
	schema := NewKeySchema(key)
	if len(schema) != 1 || *schema[0].AttributeName != "foo" {
		t.Errorf("error on NewKeySchema, actual=%v", schema)
	}

	key2 := NewKeyElement("foo2", "bar2")
	schema = NewKeySchema(key, key2)
	if len(schema) != 2 || *schema[0].AttributeName != "foo" || *schema[1].AttributeName != "foo2" {
		t.Errorf("error on NewKeySchema, actual=%v", schema)
	}
}

func TestNewKeyElement(t *testing.T) {
	key := NewKeyElement("foo", "bar")
	if *key.AttributeName != "foo" || *key.KeyType != "bar" {
		t.Errorf("error on NewKeyElement, actual=%v", key)
	}
}

func TestNewHashKeyElement(t *testing.T) {
	key := NewHashKeyElement("foo")
	if *key.AttributeName != "foo" || *key.KeyType != KeyTypeHash {
		t.Errorf("error on NewHashKeyElement, actual=%v", key)
	}
}

func TestNewRangeKeyElement(t *testing.T) {
	key := NewRangeKeyElement("foo")
	if *key.AttributeName != "foo" || *key.KeyType != KeyTypeRange {
		t.Errorf("error on NewRangeKeyElement, actual=%v", key)
	}
}

func TestNewAttributeDefinition(t *testing.T) {
	attr := NewAttributeDefinition("foo", "S")
	if *attr.AttributeName != "foo" || *attr.AttributeType != "S" {
		t.Errorf("error on NewAttributeDefinition, actual=%v", attr)
	}

	attr = NewAttributeDefinition("foo", "bar")
	if attr.AttributeName != nil || attr.AttributeType != nil {
		t.Errorf("error on NewAttributeDefinition, attributes must be nil, actual=%v", attr)
	}
}

func TestNewStringAttribute(t *testing.T) {
	attr := NewStringAttribute("foo")
	if *attr.AttributeName != "foo" || *attr.AttributeType != "S" {
		t.Errorf("error on NewStringAttribute, actual=%v", attr)
	}
}

func TestNewNumberAttribute(t *testing.T) {
	attr := NewNumberAttribute("foo")
	if *attr.AttributeName != "foo" || *attr.AttributeType != "N" {
		t.Errorf("error on NewNumberAttribute, actual=%v", attr)
	}
}

func TestNewByteAttribute(t *testing.T) {
	attr := NewByteAttribute("foo")
	if *attr.AttributeName != "foo" || *attr.AttributeType != "B" {
		t.Errorf("error on NewByteAttribute, actual=%v", attr)
	}
}

func TestNewBoolAttribute(t *testing.T) {
	attr := NewBoolAttribute("foo")
	if *attr.AttributeName != "foo" || *attr.AttributeType != "BOOL" {
		t.Errorf("error on NewBoolAttribute, actual=%v", attr)
	}
}

type TestStruct struct{}
