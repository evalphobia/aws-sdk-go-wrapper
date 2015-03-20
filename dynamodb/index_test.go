package dynamodb

import (
	"testing"

	SDK "github.com/awslabs/aws-sdk-go/gen/dynamodb"
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
}

func TestNewLSI(t *testing.T) {
	keys := NewKeySchema(NewHashKeyElement("foo"))
	lsi := NewLSI("name", keys)
	if *lsi.IndexName != "name" || *lsi.Projection.ProjectionType != SDK.ProjectionTypeAll {
		t.Errorf("error on NewLSI, actual=%v", lsi)
	}
	k := lsi.KeySchema
	if len(k) != 1 || *k[0].AttributeName != "foo" {
		t.Errorf("error on NewLSI, actual=%v", lsi)
	}
}

func TestNewGSI(t *testing.T) {
	keys := NewKeySchema(NewHashKeyElement("foo"))
	gsi := NewGSI("name", keys, NewProvisionedThroughput(5, 8))
	if *gsi.IndexName != "name" || *gsi.Projection.ProjectionType != SDK.ProjectionTypeAll {
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
}
