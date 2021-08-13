package s3

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testS3Path         = "test_path"
	testPutBucketName  = "test-put-bucket"
	testCopyBucketName = "test-copy-bucket"
	testBaseURL        = "http://localhost:4567/" + testPutBucketName
)

func testPutObject(t *testing.T) {
	assert := assert.New(t)
	createBucket(testPutBucketName)
	f := openFile(t)
	defer f.Close() // nolint:gosec

	svc := getTestClient(t)
	b, err := svc.GetBucket(testPutBucketName)
	assert.NoError(err)

	obj := NewPutObject(f)
	b.PutOne(obj, testS3Path, ACLPublicRead)

	err = b.PutAll()
	assert.NoError(err)
}

func TestAddObject(t *testing.T) {
	assert := assert.New(t)
	createBucket(testPutBucketName)
	f := openFile(t)
	defer f.Close() // nolint:gosec

	svc := getTestClient(t)
	b, err := svc.GetBucket(testPutBucketName)
	assert.NoError(err)

	obj := NewPutObject(f)
	b.AddObject(obj, testS3Path)
	assert.Equal(1, len(b.putSpool))

	obj2 := NewPutObjectString("testString")
	b.AddObject(obj2, testS3Path)
	assert.Equal(2, len(b.putSpool))

	for i, o := range []*PutObject{obj, obj2} {
		req := b.putSpool[i]
		assert.Equal("public-read", *req.ACL)
		assert.Equal(b.name, *req.Bucket)
		assert.Equal(o.data, req.Body)
		assert.Equal(o.Size(), *req.ContentLength)
		assert.Equal(o.FileType(), *req.ContentType)
		assert.Equal(testS3Path, *req.Key)
	}
}

func TestAddSecretObject(t *testing.T) {
	assert := assert.New(t)
	createBucket(testPutBucketName)
	f := openFile(t)
	defer f.Close() // nolint:gosec

	svc := getTestClient(t)
	b, err := svc.GetBucket(testPutBucketName)
	assert.NoError(err)

	obj := NewPutObject(f)
	b.AddSecretObject(obj, testS3Path)
	assert.Equal(1, len(b.putSpool))

	obj2 := NewPutObjectString("testString")
	b.AddSecretObject(obj2, testS3Path)
	assert.Equal(2, len(b.putSpool))

	for i, o := range []*PutObject{obj, obj2} {
		req := b.putSpool[i]
		assert.Equal("authenticated-read", *req.ACL)
		assert.Equal(b.name, *req.Bucket)
		assert.Equal(o.data, req.Body)
		assert.Equal(o.Size(), *req.ContentLength)
		assert.Equal(o.FileType(), *req.ContentType)
		assert.Equal(testS3Path, *req.Key)
	}
}

func TestAddPrivateObject(t *testing.T) {
	assert := assert.New(t)
	createBucket(testPutBucketName)
	f := openFile(t)
	defer f.Close() // nolint:gosec

	svc := getTestClient(t)
	b, err := svc.GetBucket(testPutBucketName)
	assert.NoError(err)

	obj := NewPutObject(f)
	b.AddPrivateObject(obj, testS3Path)
	assert.Equal(1, len(b.putSpool))

	obj2 := NewPutObjectString("testString")
	b.AddPrivateObject(obj2, testS3Path)
	assert.Equal(2, len(b.putSpool))

	for i, o := range []*PutObject{obj, obj2} {
		req := b.putSpool[i]
		assert.Equal("private", *req.ACL)
		assert.Equal(b.name, *req.Bucket)
		assert.Equal(o.data, req.Body)
		assert.Equal(o.Size(), *req.ContentLength)
		assert.Equal(o.FileType(), *req.ContentType)
		assert.Equal(testS3Path, *req.Key)
	}
}

func TestPutAll(t *testing.T) {
	assert := assert.New(t)
	createBucket(testPutBucketName)

	svc := getTestClient(t)
	b, err := svc.GetBucket(testPutBucketName)
	assert.NoError(err)

	// add spool
	obj := NewPutObjectString("testString-01")
	b.AddObject(obj, testS3Path+"_string1")
	obj2 := NewPutObjectString("testString-02")
	b.AddObject(obj2, testS3Path+"_string2")

	// write data
	err = b.PutAll()
	assert.NoError(err)

	// verify
	data, err := b.GetObjectByte(testS3Path + "_string1")
	assert.NoError(err)
	assert.Equal("testString-01", string(data))

	data, err = b.GetObjectByte(testS3Path + "_string2")
	assert.NoError(err)
	assert.Equal("testString-02", string(data))

	// Data copy error is occurred on Travis CI, Skip it.
	// f := openFile(t)
	// defer f.Close()
	// obj := NewPutObjectCopy(f)
	// b.AddObject(obj, testS3Path+"_file")
	// data, err := b.GetObjectByte(testS3Path + "_file")
	// assert.NoError(err)
	// assert.Equal(len(obj.dataByte), len(data))
}

func TestGetObjectByte(t *testing.T) {
	assert := assert.New(t)
	createBucket(testPutBucketName)
	testPutObject(t)

	f := openFile(t)
	fs, _ := f.Stat()
	defer f.Close() // nolint:gosec

	svc := getTestClient(t)
	b, err := svc.GetBucket(testPutBucketName)
	assert.NoError(err)

	// get existed data
	data, err := b.GetObjectByte(testS3Path)
	assert.NoError(err)
	assert.Equal(int(fs.Size()), len(data))

	// get from non existed path
	data, err = b.GetObjectByte("/non_exist/path")
	assert.Error(err)
	assert.Nil(data)
}

func TestGetURL(t *testing.T) {
	assert := assert.New(t)
	createBucket(testPutBucketName)

	svc := getTestClient(t)
	b, err := svc.GetBucket(testPutBucketName)
	assert.NoError(err)

	// get existed data
	url := b.GetURL(testS3Path)
	assert.Equal(url, testBaseURL+"/test_path")

	// get from non existed path
	url = b.GetURL("non_exist/path")
	assert.Equal(url, testBaseURL+"/non_exist/path")
}

func TestGetSecretURL(t *testing.T) {
	assert := assert.New(t)
	createBucket(testPutBucketName)

	svc := getTestClient(t)
	b, err := svc.GetBucket(testPutBucketName)
	assert.NoError(err)

	// get existed data
	data, err := b.GetSecretURL(testS3Path)
	assert.NoError(err)
	assert.Contains(data, testBaseURL+"/test_path?")
	assert.Contains(data, "X-Amz-Expires=180")

	// get from non existed path
	data, err = b.GetSecretURL("non_exist/path")
	assert.NoError(err)
	assert.Contains(data, testBaseURL+"/non_exist/path")
	assert.Contains(data, "X-Amz-Expires=180")
}

func TestGetSecretURLWithExpire(t *testing.T) {
	assert := assert.New(t)
	createBucket(testPutBucketName)

	svc := getTestClient(t)
	b, err := svc.GetBucket(testPutBucketName)
	assert.NoError(err)

	// get existed data
	data, err := b.GetSecretURLWithExpire(testS3Path, 520)
	assert.NoError(err)
	assert.Contains(data, testBaseURL+"/test_path?")
	assert.Contains(data, "X-Amz-Expires=520")

	// get from non existed path
	data, err = b.GetSecretURLWithExpire("non_exist/path", 10)
	assert.NoError(err)
	assert.Contains(data, testBaseURL+"/non_exist/path")
	assert.Contains(data, "X-Amz-Expires=10")
}

func TestListAllObjects(t *testing.T) {
	a := assert.New(t)
	createBucket(testPutBucketName)
	testPutObject(t)

	const (
		targetPrefix = "list-path/"
		objectCount  = 5
	)

	f := openFile(t)
	defer f.Close() // nolint:gosec
	obj := NewPutObject(f)

	svc := getTestClient(t)
	b, err := svc.GetBucket(testPutBucketName)
	a.NoError(err)

	for i := 0; i < objectCount; i++ {
		b.DeleteObject(fmt.Sprintf("%s/%d", targetPrefix, i))
	}
	list, err := b.ListAllObjects(targetPrefix)
	a.NoError(err)
	a.Len(list, 0)

	for i := 0; i < objectCount; i++ {
		b.AddObject(obj, fmt.Sprintf("%s/%d", targetPrefix, i))
	}
	err = b.PutAll()
	a.NoError(err)

	list, err = b.ListAllObjects(targetPrefix)
	a.NoError(err)
	a.Len(list, objectCount)
}

func TestCopyFrom(t *testing.T) {
	a := assert.New(t)
	createBucket(testPutBucketName)
	createBucket(testCopyBucketName)
	testPutObject(t)

	f := openFile(t)
	fs, _ := f.Stat()
	defer f.Close() // nolint:gosec

	svc := getTestClient(t)
	b1, err := svc.GetBucket(testPutBucketName)
	a.NoError(err)

	b2, err := svc.GetBucket(testCopyBucketName)
	a.NoError(err)

	// copy data
	resp, err := b2.CopyFrom(testPutBucketName, testS3Path, testS3Path)
	a.NoError(err)
	a.NotEmpty(resp.ETag)

	// check copied content
	data, err := b1.GetObjectByte(testS3Path)
	a.NoError(err)
	a.Equal(int(fs.Size()), len(data))
}

func TestCopyTo(t *testing.T) {
	a := assert.New(t)
	createBucket(testPutBucketName)
	createBucket(testCopyBucketName)
	testPutObject(t)

	f := openFile(t)
	fs, _ := f.Stat()
	defer f.Close() // nolint:gosec

	svc := getTestClient(t)
	b1, err := svc.GetBucket(testPutBucketName)
	a.NoError(err)

	b2, err := svc.GetBucket(testCopyBucketName)
	a.NoError(err)

	// copy data
	resp, err := b1.CopyTo(testS3Path, testCopyBucketName, testS3Path)
	a.NoError(err)
	a.NotEmpty(resp.ETag)

	// check copied content
	data, err := b2.GetObjectByte(testS3Path)
	a.NoError(err)
	a.Equal(int(fs.Size()), len(data))
}

func TestDeleteObject(t *testing.T) {
	assert := assert.New(t)
	createBucket(testPutBucketName)
	testPutObject(t)

	svc := getTestClient(t)
	b, err := svc.GetBucket(testPutBucketName)
	assert.NoError(err)

	// existed path
	_, errBefore := b.GetObjectByte(testS3Path)
	err = b.DeleteObject(testS3Path)
	_, errAfter := b.GetObjectByte(testS3Path)

	assert.NoError(errBefore)
	assert.NoError(err)
	assert.Error(errAfter)

	//  non existed path
	_, errBefore = b.GetObjectByte("/non_exist/path")
	err = b.DeleteObject("/non_exist/path")
	_, errAfter = b.GetObjectByte("/non_exist/path")

	assert.Error(errBefore)
	assert.NoError(err)
	assert.Error(errAfter)
}
