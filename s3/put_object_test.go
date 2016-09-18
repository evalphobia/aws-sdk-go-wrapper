package s3

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testFileName = "client_test.go"
var testFileType = "application/octet-stream"

func openFile(t *testing.T) *os.File {
	f, err := os.Open(testFileName)
	assert.NoError(t, err)
	return f
}

func TestNewPutObject(t *testing.T) {
	assert := assert.New(t)
	f := openFile(t)
	defer f.Close()

	obj := NewPutObject(f)
	assert.Equal(testFileType, obj.dataType)

	stat, _ := f.Stat()
	assert.True(obj.size > 1)
	assert.Equal(stat.Size(), obj.size)
}

func TestNewPutObjectCopy(t *testing.T) {
	assert := assert.New(t)
	f := openFile(t)
	defer f.Close()

	obj := NewPutObjectCopy(f)
	assert.Equal(testFileType, obj.dataType)

	stat, _ := f.Stat()
	assert.True(obj.size > 1)
	assert.Equal(stat.Size(), obj.size)
}

func TestNewPutObjectString(t *testing.T) {
	assert := assert.New(t)
	f := openFile(t)
	defer f.Close()

	d, err := ioutil.ReadFile(testFileName)
	assert.NoError(err)
	data := string(d)

	obj := NewPutObjectString(data)
	assert.Equal(testFileType, obj.dataType)

	stat, _ := f.Stat()
	assert.True(obj.size > 1)
	assert.Equal(stat.Size(), obj.size)
}

func TestContent(t *testing.T) {
	assert := assert.New(t)
	f := openFile(t)
	defer f.Close()

	obj := NewPutObject(f)
	assert.Equal(f, obj.Content())
}

func TestSize(t *testing.T) {
	assert := assert.New(t)
	f := openFile(t)
	defer f.Close()

	obj := NewPutObject(f)
	stat, _ := f.Stat()
	assert.Equal(stat.Size(), obj.Size())
}

func TestFileType(t *testing.T) {
	assert := assert.New(t)
	f := openFile(t)
	defer f.Close()

	obj := NewPutObject(f)
	assert.Equal(testFileType, obj.FileType())
}

func TestString(t *testing.T) {
	assert := assert.New(t)
	f := openFile(t)
	defer f.Close()

	d, err := ioutil.ReadFile(testFileName)
	assert.NoError(err)
	data := string(d)

	obj := NewPutObjectCopy(f)
	assert.Equal(data, obj.String())
}

func TestSetTypeAsText(t *testing.T) {
	assert := assert.New(t)
	f := openFile(t)
	defer f.Close()

	obj := NewPutObjectCopy(f)
	assert.Equal(testFileType, obj.FileType())

	obj.SetTypeAsText()
	assert.Equal("text/plain", obj.FileType())
}
