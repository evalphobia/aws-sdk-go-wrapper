package s3

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = fmt.Sprint("")

var testFileName = "client_test.go"
var testFileType = "application/octet-stream"

func openFile(t *testing.T) *os.File {
	f, err := os.Open(testFileName)
	assert.Nil(t, err)
	return f
}

func TestNewS3Object(t *testing.T) {
	f := openFile(t)
	defer f.Close()

	obj := NewS3Object(f)
	assert.Equal(t, testFileType, obj.dataType)

	stat, _ := f.Stat()
	assert.True(t, obj.size > 1)
	assert.Equal(t, stat.Size(), obj.size)
}

func TestNewS3ObjectCopy(t *testing.T) {
	f := openFile(t)
	defer f.Close()

	obj := NewS3ObjectCopy(f)
	assert.Equal(t, testFileType, obj.dataType)

	stat, _ := f.Stat()
	assert.True(t, obj.size > 1)
	assert.Equal(t, stat.Size(), obj.size)
}

func TestNewS3ObjectString(t *testing.T) {
	f := openFile(t)
	defer f.Close()
	d, err := ioutil.ReadFile(testFileName)
	assert.Nil(t, err)
	data := string(d)

	obj := NewS3ObjectString(&data)
	assert.Equal(t, testFileType, obj.dataType)

	stat, _ := f.Stat()
	assert.True(t, obj.size > 1)
	assert.Equal(t, stat.Size(), obj.size)
}

func TestContent(t *testing.T) {
	f := openFile(t)
	defer f.Close()

	obj := NewS3Object(f)
	assert.Equal(t, f, obj.Content())
}

func TestSize(t *testing.T) {
	f := openFile(t)
	defer f.Close()

	obj := NewS3Object(f)
	stat, _ := f.Stat()
	assert.Equal(t, stat.Size(), obj.Size())
}

func TestFileType(t *testing.T) {
	f := openFile(t)
	defer f.Close()

	obj := NewS3Object(f)
	assert.Equal(t, testFileType, obj.FileType())
}

func TestString(t *testing.T) {
	f := openFile(t)
	defer f.Close()
	d, err := ioutil.ReadFile(testFileName)
	assert.Nil(t, err)
	data := string(d)

	obj := NewS3ObjectCopy(f)
	assert.Equal(t, data, obj.String())
}

func TestSetTypeAsText(t *testing.T) {
	f := openFile(t)
	defer f.Close()

	obj := NewS3ObjectCopy(f)
	assert.Equal(t, testFileType, obj.FileType())

	obj.SetTypeAsText()
	assert.Equal(t, "text/plain", obj.FileType())
}

func TestNewFileBuffer(t *testing.T) {
	var b []byte

	buf := NewFileBuffer(b)
	assert.Equal(t, int64(0), buf.Index)
}

func TestClose(t *testing.T) {
	var b []byte
	buf := NewFileBuffer(b)
	err := buf.Close()
	assert.Nil(t, err)
}

func TestBytes(t *testing.T) {
	str := "abcd"
	b := []byte(str)
	buf := NewFileBuffer(b)
	assert.Equal(t, string(b), string(buf.Bytes()))
}

func TestRead(t *testing.T) {
	str := "abcd"
	b := []byte(str)
	buf := NewFileBuffer(b)

	n, err := buf.Read([]byte("a"))
	assert.Nil(t, err)
	assert.Equal(t, 1, n)

	n, err = buf.Read(b)
	assert.Nil(t, err)
	assert.Equal(t, 3, n)

	n, err = buf.Read(b)
	assert.NotNil(t, err)
	assert.Equal(t, "EOF", err.Error())
	assert.Equal(t, 0, n)
}

func TestSeek(t *testing.T) {
	str := "abcd"
	b := []byte(str)
	buf := NewFileBuffer(b)

	n, err := buf.Seek(2, 0)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), n)

	n, err = buf.Seek(2, 1)
	assert.NotNil(t, err)
	assert.Equal(t, "Unsupported Seek Method.", err.Error())
	assert.Equal(t, int64(0), n)

	n, err = buf.Seek(9, 0)
	assert.NotNil(t, err)
	assert.Equal(t, "Seek: Invalid Offset", err.Error())
	assert.Equal(t, int64(0), n)
}
