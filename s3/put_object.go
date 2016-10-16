package s3

import (
	"bytes"
	"io"
	"os"
)

const (
	mimeBinary = "application/octet-stream"
	mimeText   = "text/plain"
)

// PutObject is wrapper struct for the object to upload S3.
type PutObject struct {
	data     io.ReadSeeker
	dataType string
	dataByte []byte
	size     int64
}

// Create new PutObject
func newPutObject(f io.ReadSeeker, size int64, typ string) *PutObject {
	return &PutObject{
		data:     f,
		dataType: typ,
		size:     size,
	}
}

// NewPutObject returns initialized *PutObject from File.
func NewPutObject(file *os.File) *PutObject {
	fi, _ := file.Stat()
	return newPutObject(file, fi.Size(), mimeBinary)
}

// NewPutObjectCopy returns initialized *PutObject from File and copy byte data.
func NewPutObjectCopy(file *os.File) *PutObject {
	buf := new(bytes.Buffer)
	io.Copy(buf, file)
	fi, _ := file.Stat()
	o := newPutObject(file, fi.Size(), mimeBinary)
	o.dataByte = buf.Bytes()
	return o
}

// NewPutObjectString returns initialized *PutObject from string.
func NewPutObjectString(s string) *PutObject {
	b := []byte(s)
	o := newPutObject(bytes.NewReader(b), int64(len(b)), mimeBinary)
	o.dataByte = b
	return o
}

func (o *PutObject) String() string {
	return string(o.dataByte)
}

// Content returns the content of the Object.
func (o *PutObject) Content() io.ReadSeeker {
	return o.data
}

// Size returns size of the content.
func (o *PutObject) Size() int64 {
	return o.size
}

// FileType returns file type of the content.
func (o *PutObject) FileType() string {
	return o.dataType
}

// SetTypeAsText sets MIME type as text file.
func (o *PutObject) SetTypeAsText() {
	o.dataType = mimeText
}
