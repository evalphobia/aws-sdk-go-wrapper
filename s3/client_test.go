package s3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	s := NewClient()
	assert.NotNil(t, s.client)

	c := s.client
	assert.Equal(t, "s3", c.ServiceName)
	assert.Equal(t, defaultEndpoint, c.Endpoint)
}

func TestGetBucket(t *testing.T) {
	s := NewClient()
	b := s.GetBucket("test")
	bucketName := defaultBucketPrefix + "test"

	assert.Equal(t, bucketName, b.name)
	assert.NotNil(t, s.buckets[bucketName])
	assert.Equal(t, b, s.buckets[bucketName])
	assert.Equal(t, b.client, s.client)

	b2 := s.GetBucket("test")
	assert.Equal(t, b, b2)
}
