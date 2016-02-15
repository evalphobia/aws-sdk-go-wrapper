package s3

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "github.com/evalphobia/aws-sdk-go-wrapper/config/json"
)

func setTestEnv() {
	os.Clearenv()
	os.Setenv("AWS_ACCESS_KEY_ID", "access")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
}

func TestNewClient(t *testing.T) {
	setTestEnv()

	s := NewClient()
	assert.NotNil(t, s.client)

	c := s.client
	assert.Equal(t, "s3", c.ServiceName)
	assert.Equal(t, defaultEndpoint, c.Endpoint)
}

func TestGetBucket(t *testing.T) {
	setTestEnv()

	s := NewClient()
	b := s.GetBucket("test")
	bucketName := defaultBucketPrefix + "test"

	assert.Equal(t, bucketName, b.name)
	assert.NotNil(t, s.buckets[bucketName])
	assert.Equal(t, b, s.buckets[bucketName])
	assert.Equal(t, b.client, s.client)

	b2 := s.GetBucket("test")
	assert.Equal(t, b, b2)

	s2 := NewClient()
	s2.SetBucketPrefix("ass-")
	b3 := s2.GetBucket("test")
	assert.Equal(t, "ass-test", b3.name)
}
