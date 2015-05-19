package dynamodb

import (
	"testing"
)

func TestNewDynamoIndex(t *testing.T) {
	hashkey := NewHashKeyElement("foo")
	idx := NewDynamoIndex("name", "local", NewKeySchema(hashkey))
	keys := idx.KeySchema
	if idx.Name != "name" || idx.IndexType != "local" || len(keys) != 1 {
		t.Errorf("error on NewDynamoIndex, actual=%v", idx)
	}
	if *keys[0].AttributeName != "foo" {
		t.Errorf("error on NewDynamoIndex.KeySchema, actual=%v", keys)
	}

	rangekey := NewRangeKeyElement("bar")
	idx2 := NewDynamoIndex("name", "local", NewKeySchema(hashkey, rangekey))
	keys = idx2.KeySchema
	if idx.Name != "name" || idx.IndexType != "local" || len(keys) != 2 {
		t.Errorf("error on NewDynamoIndex, actual=%v", idx)
	}
	if *keys[0].AttributeName != "foo" || *keys[1].AttributeName != "bar" {
		t.Errorf("error on NewDynamoIndex.KeySchema, actual=%v", keys)
	}

}

func TestGetHashKeyName(t *testing.T) {
	hashkey := NewHashKeyElement("foo")
	rangekey := NewRangeKeyElement("bar")
	idx := NewDynamoIndex("name", "local", NewKeySchema(hashkey, rangekey))
	hashName := idx.GetHashKeyName()
	if hashName != "foo" {
		t.Errorf("error on GetHashKeyName, actual=%v", hashName)
	}
}

func TestGetRangeKeyName(t *testing.T) {
	hashkey := NewHashKeyElement("foo")
	rangekey := NewRangeKeyElement("bar")
	idx := NewDynamoIndex("name", "local", NewKeySchema(hashkey, rangekey))
	rangeName := idx.GetRangeKeyName()
	if rangeName != "bar" {
		t.Errorf("error on GetHashKeyName, actual=%v", rangeName)
	}

	idx2 := NewDynamoIndex("name", "local", NewKeySchema(hashkey))
	rangeName = idx2.GetRangeKeyName()
	if rangeName != "" {
		t.Errorf("error on GetHashKeyName, actual=%v", rangeName)
	}
}

func TestNewLSI(t *testing.T) {
	keys := NewKeySchema(NewHashKeyElement("foo"))
	lsi := NewLSI("name", keys)
	if *lsi.IndexName != "name" || *lsi.Projection.ProjectionType != ProjectionTypeAll {
		t.Errorf("error on NewLSI, actual=%v", lsi)
	}
	k := lsi.KeySchema
	if len(k) != 1 || *k[0].AttributeName != "foo" {
		t.Errorf("error on NewLSI, actual=%v", lsi)
	}

	lsi2 := NewLSI("name", keys, ProjectionTypeKeysOnly)
	if *lsi2.IndexName != "name" || *lsi2.Projection.ProjectionType != ProjectionTypeKeysOnly {
		t.Errorf("error on NewLSI, actual=%v", lsi2)
	}
	k = lsi2.KeySchema
	if len(k) != 1 || *k[0].AttributeName != "foo" {
		t.Errorf("error on NewLSI, actual=%v", lsi2)
	}
}

func TestNewGSI(t *testing.T) {
	keys := NewKeySchema(NewHashKeyElement("foo"))
	gsi := NewGSI("name", keys, NewProvisionedThroughput(5, 8))
	if *gsi.IndexName != "name" || *gsi.Projection.ProjectionType != ProjectionTypeAll {
		t.Errorf("error on NewGSI, actual=%v", gsi)
	}
	k := gsi.KeySchema
	if len(k) != 1 || *k[0].AttributeName != "foo" {
		t.Errorf("error on NewGSI, actual=%v", gsi)
	}
	tp := gsi.ProvisionedThroughput
	if *tp.ReadCapacityUnits != 5 || *tp.WriteCapacityUnits != 8 {
		t.Errorf("error on NewGSI, actual=%v", gsi)
	}

	gsi2 := NewGSI("name", keys, NewProvisionedThroughput(5, 8), ProjectionTypeKeysOnly)
	if *gsi2.IndexName != "name" || *gsi2.Projection.ProjectionType != ProjectionTypeKeysOnly {
		t.Errorf("error on NewGSI, actual=%v", gsi)
	}
	k = gsi2.KeySchema
	if len(k) != 1 || *k[0].AttributeName != "foo" {
		t.Errorf("error on NewGSI, actual=%v", gsi2)
	}
	tp = gsi2.ProvisionedThroughput
	if *tp.ReadCapacityUnits != 5 || *tp.WriteCapacityUnits != 8 {
		t.Errorf("error on NewGSI, actual=%v", gsi2)
	}
}
