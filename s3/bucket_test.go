package s3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testS3Path = "/test_path"
var testBucketName = "test-bucket"

func TestAddObject(t *testing.T) {
	assert := assert.New(t)
	setTestEnv()
	f := openFile(t)
	defer f.Close()
	obj := NewS3Object(f)

	s := NewClient()
	b := s.GetBucket(testBucketName)
	b.AddObject(obj, testS3Path)

	assert.Equal(1, len(b.objects))

	req := b.objects[0]
	assert.Equal("public-read", *req.ACL)
	assert.Equal(b.name, *req.Bucket)
	assert.Equal(obj.data, req.Body)
	assert.Equal(obj.Size(), *req.ContentLength)
	assert.Equal(obj.FileType(), *req.ContentType)
	assert.Equal(testS3Path, *req.Key)
}

func TestAddSecretObject(t *testing.T) {
	assert := assert.New(t)
	setTestEnv()
	f := openFile(t)
	defer f.Close()
	obj := NewS3Object(f)

	s := NewClient()
	b := s.GetBucket(testBucketName)
	b.AddSecretObject(obj, testS3Path)

	assert.Equal(1, len(b.objects))

	req := b.objects[0]
	assert.Equal("authenticated-read", *req.ACL)
	assert.Equal(b.name, *req.Bucket)
	assert.Equal(obj.data, req.Body)
	assert.Equal(obj.Size(), *req.ContentLength)
	assert.Equal(obj.FileType(), *req.ContentType)
	assert.Equal(testS3Path, *req.Key)
}

func TestPut(t *testing.T) {
	assert := assert.New(t)
	setTestEnv()
	f := openFile(t)
	defer f.Close()
	obj := NewS3Object(f)

	s := NewClient()
	b := s.GetBucket(testBucketName)
	b.AddObject(obj, testS3Path)

	err := b.Put()
	assert.Nil(err)
}

func TestGetObjectByte(t *testing.T) {
	assert := assert.New(t)
	setTestEnv()
	TestPut(t)

	f := openFile(t)
	fs, _ := f.Stat()
	defer f.Close()

	s := NewClient()
	b := s.GetBucket(testBucketName)

	// get existed data
	data, err := b.GetObjectByte(testS3Path)
	assert.Nil(err)
	assert.Equal(int(fs.Size()), len(data))

	// get from non existed path
	data, err = b.GetObjectByte("/non_exist/path")
	assert.NotNil(err)
	assert.Equal([]byte{}, data)
}

func TestGetURL(t *testing.T) {
	assert := assert.New(t)
	setTestEnv()
	s := NewClient()
	b := s.GetBucket(testBucketName)

	baseURL := "http://localhost:4567/" + defaultBucketPrefix + testBucketName

	// get existed data
	url := b.GetURL(testS3Path)
	assert.Equal(url, baseURL+"/test_path")

	// get from non existed path
	url = b.GetURL("/non_exist/path")
	assert.Equal(url, baseURL+"/non_exist/path")
}

func TestGetSecretURL(t *testing.T) {
	assert := assert.New(t)
	setTestEnv()
	s := NewClient()
	b := s.GetBucket(testBucketName)

	baseURL := "http://localhost:4567/" + defaultBucketPrefix + testBucketName

	// get existed data
	data, err := b.GetSecretURL(testS3Path)
	assert.Nil(err)
	assert.Contains(data, baseURL+"/test_path?")
	assert.Contains(data, "X-Amz-Expires=180")

	// get from non existed path
	data, err = b.GetSecretURL("/non_exist/path")
	assert.Nil(err)
	assert.Contains(data, baseURL+"/non_exist/path")
	assert.Contains(data, "X-Amz-Expires=180")
}

func TestGetSecretURLWithExpire(t *testing.T) {
	assert := assert.New(t)
	setTestEnv()
	s := NewClient()
	b := s.GetBucket(testBucketName)

	baseURL := "http://localhost:4567/" + defaultBucketPrefix + testBucketName

	// get existed data
	data, err := b.GetSecretURLWithExpire(testS3Path, 520)
	assert.Nil(err)
	assert.Contains(data, baseURL+"/test_path?")
	assert.Contains(data, "X-Amz-Expires=520")

	// get from non existed path
	data, err = b.GetSecretURLWithExpire("/non_exist/path", 10)
	assert.Nil(err)
	assert.Contains(data, baseURL+"/non_exist/path")
	assert.Contains(data, "X-Amz-Expires=10")
}

func TestDeleteObject(t *testing.T) {
	assert := assert.New(t)
	setTestEnv()
	TestPut(t)

	s := NewClient()
	b := s.GetBucket(testBucketName)

	// existed path
	_, errBefore := b.GetObjectByte(testS3Path)
	err := b.DeleteObject(testS3Path)
	_, errAfter := b.GetObjectByte(testS3Path)

	assert.Nil(errBefore)
	assert.Nil(err)
	assert.NotNil(errAfter)

	//  non existed path
	_, errBefore = b.GetObjectByte("/non_exist/path")
	err = b.DeleteObject("/non_exist/path")
	_, errAfter = b.GetObjectByte("/non_exist/path")

	assert.NotNil(errBefore)
	assert.Nil(err)
	assert.NotNil(errAfter)
}
