package kinesis

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewStream(t *testing.T) {
	assert := assert.New(t)

	recreateTestStream(t)

	tests := []struct {
		name      string
		isSuccess bool
	}{
		{"foo", false},
		{testStreamName, true},
	}

	svc := getTestClient(t)
	for _, tt := range tests {
		target := fmt.Sprintf("%+v", tt)

		s, err := NewStream(svc, tt.name)
		if !tt.isSuccess {
			assert.Error(err, target)
			assert.Nil(s, target)
			continue
		}

		assert.NoError(err, target)
		assert.NotNil(s, target)
		assert.Equal(tt.name, s.nameWithPrefix, target)
		assert.True(len(s.shardIDs) > 0, target, s.shardIDs)
	}
}

func TestGetShardIDs(t *testing.T) {
	assert := assert.New(t)

	testStreamName2 := "test-stream-2"
	svc := getTestClient(t)
	svc.CreateStreamWithName(testStreamName2)

	recreateTestStream(t)

	tests := []struct {
		name      string
		isSuccess bool
	}{
		{"foo", false},
		{"bar", false},
		{testStreamName, true},
		{testStreamName2, true},
	}

	for _, tt := range tests {
		target := fmt.Sprintf("%+v", tt)

		s := &Stream{
			service:        svc,
			nameWithPrefix: tt.name,
		}
		ids, err := s.GetShardIDs()

		if !tt.isSuccess {
			assert.Error(err, target)
			assert.Nil(ids, target)
			continue
		}
		assert.NoError(err, target)
		assert.True(len(ids) > 0, target, ids)
	}
}

func TestGetLatestRecords(t *testing.T) {
	t.Skip("unstable test")

	assert := assert.New(t)
	recreateTestStream(t)
	s := getTestStream(t)

	// put data during GetLatestRecords... :(
	go func() {
		s.PutRecord([]byte("TestGetLatestRecords"))
		s.PutRecord([]byte("TestGetLatestRecords"))
		s.PutRecord([]byte("TestGetLatestRecords"))
		s.PutRecord([]byte("TestGetLatestRecords"))
		s.PutRecord([]byte("TestGetLatestRecords"))
		s.PutRecord([]byte("TestGetLatestRecords"))
		s.PutRecord([]byte("TestGetLatestRecords"))
		s.PutRecord([]byte("TestGetLatestRecords"))
		s.PutRecord([]byte("TestGetLatestRecords"))
	}()
	time.Sleep(1 * time.Millisecond)

	list, err := s.GetLatestRecords()
	assert.NoError(err)
	assert.Len(list, 1)

	res := list[0]
	assert.True(res.Count > 0)
	assert.NotEmpty(res.NextShardIterator)
	assert.True(len(res.Items) > 0)
	assert.Equal("TestGetLatestRecords", string(res.Items[0].Data))
}

func TestGetRecords(t *testing.T) {
	assert := assert.New(t)
	recreateTestStream(t)
	s := getTestStream(t)

	// empty result
	s.PutRecord([]byte("TestGetRecords1"))
	result, err := s.GetRecords(GetCondition{
		ShardID:           s.shardIDs[0],
		ShardIteratorType: IteratorTypeLatest,
	})
	assert.NoError(err)
	assert.Equal(0, result.Count)
	assert.Equal(0, len(result.Items))

	// get 1 record from the previous
	s.PutRecord([]byte("TestGetRecords2"))
	result, err = s.GetRecords(GetCondition{
		ShardID:           s.shardIDs[0],
		ShardIteratorType: IteratorTypeLatest,
		ShardIterator:     result.NextShardIterator,
	})
	assert.NoError(err)
	assert.Equal(1, result.Count)
	assert.Equal(1, len(result.Items))
	assert.Equal("TestGetRecords2", string(result.Items[0].Data))

	// get 2 record from the beginning
	result, err = s.GetRecords(GetCondition{
		ShardID:           s.shardIDs[0],
		ShardIteratorType: IteratorTypeTrimHorizon,
	})
	assert.NoError(err)
	assert.Equal(2, result.Count)
	assert.Equal(2, len(result.Items))
	assert.Equal("TestGetRecords1", string(result.Items[0].Data))
	assert.Equal("TestGetRecords2", string(result.Items[1].Data))
}

func TestPutRecord(t *testing.T) {
	assert := assert.New(t)
	recreateTestStream(t)
	s := getTestStream(t)

	// before put
	result, err := s.GetRecords(GetCondition{
		ShardID:           s.shardIDs[0],
		ShardIteratorType: IteratorTypeTrimHorizon,
	})
	assert.NoError(err)
	assert.Equal(0, result.Count)
	assert.Equal(0, len(result.Items))

	// execute put
	err = s.PutRecord([]byte("TestPutRecord"))
	assert.NoError(err)

	// after put
	result, err = s.GetRecords(GetCondition{
		ShardID:           s.shardIDs[0],
		ShardIteratorType: IteratorTypeTrimHorizon,
	})
	assert.NoError(err)
	assert.Equal(1, result.Count)
	assert.Equal(1, len(result.Items))
	assert.Equal("TestPutRecord", string(result.Items[0].Data))
}
