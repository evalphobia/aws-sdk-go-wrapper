package s3

import (
	"testing"

	SDK "github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
)

const (
	defaultEndpoint     = "http://localhost:4567"
	testEmptyBucketName = "test-empty-bucket"
)

func getTestConfig() config.Config {
	return config.Config{
		AccessKey:        "access",
		SecretKey:        "secret",
		Endpoint:         defaultEndpoint,
		S3ForcePathStyle: true,
	}
}

func getTestClient(t *testing.T) *S3 {
	svc, err := New(getTestConfig())
	if err != nil {
		t.Errorf("error on create client; error=%s;", err.Error())
		t.FailNow()
	}
	return svc
}

func createBucket(name string) error {
	svc, err := New(getTestConfig())
	if err != nil {
		return err
	}

	ok, err := svc.IsExistBucket(name)
	switch {
	case err != nil:
		return err
	case ok:
		return nil
	}

	return svc.CreateBucketWithName(name)
}

func TestNew(t *testing.T) {
	assert := assert.New(t)

	svc, err := New(getTestConfig())
	assert.NoError(err)
	assert.NotNil(svc.client)
	assert.Equal("s3", svc.client.ServiceName)
	assert.Equal(defaultEndpoint, svc.client.Endpoint)

	region := "us-west-1"
	svc, err = New(config.Config{
		Region: region,
	})
	assert.NoError(err)
	expectedEndpoint := "https://s3." + region + ".amazonaws.com"
	assert.Equal(expectedEndpoint, svc.client.Endpoint)
}

func TestSetLogger(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)
	assert.Equal(log.DefaultLogger, svc.logger)

	stdLogger := &log.StdLogger{}
	svc.SetLogger(stdLogger)
	assert.Equal(stdLogger, svc.logger)
}

func TestGetBucket(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)
	testPrefix := "prefix-"
	createBucket(testEmptyBucketName)
	createBucket(testPrefix + testEmptyBucketName)

	b, err := svc.GetBucket(testEmptyBucketName)
	assert.NoError(err)

	assert.Equal(testEmptyBucketName, b.name)
	assert.NotNil(svc.buckets[testEmptyBucketName])
	assert.Equal(b, svc.buckets[testEmptyBucketName])
	assert.Equal(b.service, svc)

	b2, err := svc.GetBucket(testEmptyBucketName)
	assert.NoError(err)
	assert.Equal(b, b2)

	svc2 := getTestClient(t)
	svc2.prefix = testPrefix

	b3, err := svc2.GetBucket(testEmptyBucketName)
	assert.NoError(err)
	assert.Equal(testPrefix+testEmptyBucketName, b3.nameWithPrefix)
}

func TestIsExistBucket(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	svc.ForceDeleteBucket(testEmptyBucketName)
	has, err := svc.IsExistBucket(testEmptyBucketName)
	assert.NoError(err)
	assert.False(has, "should be deleted")

	createBucket(testEmptyBucketName)
	has, err = svc.IsExistBucket(testEmptyBucketName)
	assert.NoError(err)
	assert.True(has, "should be created")
}

func TestCreateBucket(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	svc.ForceDeleteBucket(testEmptyBucketName)
	has, err := svc.IsExistBucket(testEmptyBucketName)
	assert.NoError(err)
	assert.False(has)

	bucketName := testEmptyBucketName
	input := &SDK.CreateBucketInput{
		Bucket: &bucketName,
	}

	err = svc.CreateBucket(input)
	assert.NoError(err, "success creation")

	has, err = svc.IsExistBucket(testEmptyBucketName)
	assert.NoError(err)
	assert.True(has, "should be created")
}

func TestCreateBucketWithName(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	svc.ForceDeleteBucket(testEmptyBucketName)
	has, err := svc.IsExistBucket(testEmptyBucketName)
	assert.NoError(err)
	assert.False(has)

	err = svc.CreateBucketWithName(testEmptyBucketName)
	assert.NoError(err, "success creation")

	has, err = svc.IsExistBucket(testEmptyBucketName)
	assert.NoError(err)
	assert.True(has, "should be created")
}

func TestForceDeleteBucket(t *testing.T) {
	assert := assert.New(t)
	svc := getTestClient(t)

	createBucket(testEmptyBucketName)
	has, err := svc.IsExistBucket(testEmptyBucketName)
	assert.NoError(err)
	assert.True(has)

	svc.ForceDeleteBucket(testEmptyBucketName)
	has, err = svc.IsExistBucket(testEmptyBucketName)
	assert.NoError(err)
	assert.False(has, "should be deleted")

	err = svc.ForceDeleteBucket(testEmptyBucketName)
	assert.Error(err, "already deleted")
}
