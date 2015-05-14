package s3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testS3Path = "/test_path"

func TestAddObject(t *testing.T) {
	f := openFile(t)
	defer f.Close()
	obj := NewS3Object(f)

	s := NewClient()
	b := s.GetBucket("test")
	b.AddObject(obj, testS3Path)

	assert.Equal(t, 1, len(b.objects))

	req := b.objects[0]
	assert.Equal(t, "public-read", *req.ACL)
	assert.Equal(t, b.name, *req.Bucket)
	assert.Equal(t, obj.data, req.Body)
	assert.Equal(t, obj.Size(), *req.ContentLength)
	assert.Equal(t, obj.FileType(), *req.ContentType)
	assert.Equal(t, testS3Path, *req.Key)
}

func TestAddSecretObject(t *testing.T) {
	f := openFile(t)
	defer f.Close()
	obj := NewS3Object(f)

	s := NewClient()
	b := s.GetBucket("test")
	b.AddSecretObject(obj, testS3Path)

	assert.Equal(t, 1, len(b.objects))

	req := b.objects[0]
	assert.Equal(t, "authenticated-read", *req.ACL)
	assert.Equal(t, b.name, *req.Bucket)
	assert.Equal(t, obj.data, req.Body)
	assert.Equal(t, obj.Size(), *req.ContentLength)
	assert.Equal(t, obj.FileType(), *req.ContentType)
	assert.Equal(t, testS3Path, *req.Key)
}

func TestPut(t *testing.T) {
	f := openFile(t)
	defer f.Close()
	obj := NewS3Object(f)

	s := NewClient()
	b := s.GetBucket("test")
	b.AddObject(obj, testS3Path)

	err := b.Put()
	assert.Nil(t, err)
}

func TestGetObjectByte(t *testing.T) {
	TestPut(t)

	f := openFile(t)
	fs, _ := f.Stat()
	defer f.Close()

	s := NewClient()
	b := s.GetBucket("test")

	// get existed data
	data, err := b.GetObjectByte(testS3Path)
	assert.Nil(t, err)
	assert.Equal(t, int(fs.Size()), len(data))

	// get from non existed path
	data, err = b.GetObjectByte("/non_exist/path")
	assert.NotNil(t, err)
	assert.Equal(t, []byte{}, data)
}

func TestGetURL(t *testing.T) {
	s := NewClient()
	b := s.GetBucket("test")

	// get existed data
	data, err := b.GetURL(testS3Path)
	assert.Nil(t, err)
	assert.Contains(t, data, "http://localhost:4567/dev-test/test_path?")

	// get from non existed path
	data, err = b.GetURL("/non_exist/path")
	assert.Nil(t, err)
	assert.Contains(t, data, "http://localhost:4567/dev-test/non_exist/path?")
}

func TestGetSecretURL(t *testing.T) {
	s := NewClient()
	b := s.GetBucket("test")

	// get existed data
	data, err := b.GetSecretURL(testS3Path)
	assert.Nil(t, err)
	assert.Contains(t, data, "http://localhost:4567/dev-test/test_path?")
	assert.Contains(t, data, "X-Amz-Expires=180")

	// get from non existed path
	data, err = b.GetSecretURL("/non_exist/path")
	assert.Nil(t, err)
	assert.Contains(t, data, "http://localhost:4567/dev-test/non_exist/path?")
	assert.Contains(t, data, "X-Amz-Expires=180")
}

func TestGetSecretURLWithExpire(t *testing.T) {
	s := NewClient()
	b := s.GetBucket("test")

	// get existed data
	data, err := b.GetSecretURLWithExpire(testS3Path, 520)
	assert.Nil(t, err)
	assert.Contains(t, data, "http://localhost:4567/dev-test/test_path?")
	assert.Contains(t, data, "X-Amz-Expires=520")

	// get from non existed path
	data, err = b.GetSecretURLWithExpire("/non_exist/path", 10)
	assert.Nil(t, err)
	assert.Contains(t, data, "http://localhost:4567/dev-test/non_exist/path?")
	assert.Contains(t, data, "X-Amz-Expires=10")
}

func TestDeleteObject(t *testing.T) {
	TestPut(t)

	s := NewClient()
	b := s.GetBucket("test")

	// existed path
	_, errBefore := b.GetObjectByte(testS3Path)
	err := b.DeleteObject(testS3Path)
	_, errAfter := b.GetObjectByte(testS3Path)

	assert.Nil(t, errBefore)
	assert.Nil(t, err)
	assert.NotNil(t, errAfter)

	//  non existed path
	_, errBefore = b.GetObjectByte("/non_exist/path")
	err = b.DeleteObject("/non_exist/path")
	_, errAfter = b.GetObjectByte("/non_exist/path")

	assert.NotNil(t, errBefore)
	assert.Nil(t, err)
	assert.NotNil(t, errAfter)
}
