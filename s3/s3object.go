// S3 object manipuration

package s3

import (
	"github.com/evalphobia/aws-sdk-go-wrapper/log"

	"bytes"
	"errors"
	"io"
	"os"
	"unsafe"
)

const (
	mimeBinary = "application/octet-stream"
	mimeText   = "text/plain"
)

// struct for S3 upload
type S3Object struct {
	data     io.Reader
	dataType string
	dataByte []byte
	size     int64
}

// Create new S3Object
func newS3Object(f io.Reader, size int64, typ string) *S3Object {
	return &S3Object{
		data:     f,
		dataType: typ,
		size:     size,
	}
}

// Create new S3Object From File
func NewS3Object(file *os.File) *S3Object {
	fi, _ := file.Stat()
	return newS3Object(file, fi.Size(), mimeBinary)
}

// Create new S3Object From File and copy byte data
func NewS3ObjectCopy(file *os.File) *S3Object {
	buf := new(bytes.Buffer)
	io.Copy(buf, file)
	fi, _ := file.Stat()
	o := newS3Object(file, fi.Size(), mimeBinary)
	o.dataByte = buf.Bytes()
	return o
}

// Create new S3Object From string
func NewS3ObjectString(str *string) *S3Object {
	b := *(*[]byte)(unsafe.Pointer(str))
	var r io.Reader = NewFileBuffer(b)
	o := newS3Object(r, int64(len(*str)), mimeBinary)
	o.dataByte = b
	return o
}

// get content from S3Object
func (o *S3Object) Content() io.Reader {
	return o.data
}

func (o *S3Object) Size() int64 {
	return o.size
}

func (o *S3Object) FileType() string {
	return o.dataType
}

func (o *S3Object) String() string {
	return string(o.dataByte)
}

func (o *S3Object) SetTypeAsText() {
	o.dataType = mimeText
}

// wrapped struct for io.ReadCloser alternative
type FileBuffer struct {
	Buffer bytes.Buffer
	Index  int64
}

func NewFileBuffer(b []byte) *FileBuffer {
	return &FileBuffer{
		Buffer: *bytes.NewBuffer(b),
		Index:  int64(0),
	}
}

func (f *FileBuffer) Bytes() []byte {
	return f.Buffer.Bytes()
}

func (f *FileBuffer) Read(p []byte) (int, error) {
	n, err := bytes.NewBuffer(f.Buffer.Bytes()[f.Index:]).Read(p)
	if err == nil {
		if f.Index+int64(len(p)) < int64(f.Buffer.Len()) {
			f.Index += int64(len(p))
		} else {
			f.Index = int64(f.Buffer.Len())
		}
	} else {
		log.Warn("[S3Object] error on read file", err.Error())
	}
	return n, err
}

func (f *FileBuffer) Seek(offset int64, whence int) (int64, error) {
	var err error
	var Index int64 = 0
	switch whence {
	case 0:
		if offset >= int64(f.Buffer.Len()) || offset < 0 {
			err = errors.New("Seek: Invalid Offset")
		} else {
			f.Index = offset
			Index = offset
		}
	default:
		err = errors.New("Unsupported Seek Method.")
	}
	if err != nil {
		log.Warn("[S3Object] error on seek file", err.Error())
	}
	return Index, err
}
