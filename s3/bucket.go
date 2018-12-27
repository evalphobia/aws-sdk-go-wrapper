package s3

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"

	SDK "github.com/aws/aws-sdk-go/service/s3"

	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

// ACL settings
const (
	ACLAuthenticatedRead = "authenticated-read"
	ACLPrivate           = "private"
	ACLPublicRead        = "public-read"
	ACLPublicReadWrite   = "public-read-write"
)

const (
	defaultExpireSecond = 180
)

// Bucket is S3 Bucket wrapper struct.
type Bucket struct {
	service *S3

	name           string
	nameWithPrefix string
	endpoint       string
	expireSecond   int

	putSpoolMu sync.Mutex
	putSpool   []*SDK.PutObjectInput
}

// NewBucket returns initialized *Bucket.
func NewBucket(svc *S3, name string) *Bucket {
	bucketName := svc.prefix + name
	return &Bucket{
		service:        svc,
		name:           name,
		nameWithPrefix: bucketName,
		endpoint:       svc.endpoint,
		expireSecond:   defaultExpireSecond,
	}
}

// SetExpire sets default expire sec for ACL access.
func (b *Bucket) SetExpire(sec int) {
	b.expireSecond = sec
}

// AddObject adds object to write spool (w/ public read access).
func (b *Bucket) AddObject(obj *PutObject, path string) {
	b.addObject(obj, path, ACLPublicRead)
}

// AddSecretObject adds object to write spool (w/ ACL permission).
func (b *Bucket) AddSecretObject(obj *PutObject, path string) {
	b.addObject(obj, path, ACLAuthenticatedRead)
}

// addObject adds object to write spool.
func (b *Bucket) addObject(obj *PutObject, path, acl string) {
	b.putSpoolMu.Lock()
	defer b.putSpoolMu.Unlock()

	size := obj.Size()
	req := &SDK.PutObjectInput{
		ACL:           &acl,
		Bucket:        &b.nameWithPrefix,
		Body:          obj.data,
		ContentLength: &size,
		ContentType:   pointers.String(obj.FileType()),
		Key:           pointers.String(path),
	}
	b.putSpool = append(b.putSpool, req)
}

// PutAll executes PutObject operation in the put spool.
func (b *Bucket) PutAll() error {
	b.putSpoolMu.Lock()
	defer b.putSpoolMu.Unlock()

	errList := newErrors()
	cli := b.service.client
	for _, obj := range b.putSpool {
		_, err := cli.PutObject(obj)
		if err != nil {
			b.service.Errorf("error on `PutObject` operation; bucket=%s; error=%s;", b.nameWithPrefix, err.Error())
			errList.Add(err)
		}
	}
	b.putSpool = nil

	if errList.HasError() {
		return errList
	}
	return nil
}

// PutOne executes PutObject operation in the put spool.
func (b *Bucket) PutOne(obj *PutObject, path, acl string) error {
	size := obj.Size()
	req := &SDK.PutObjectInput{
		ACL:           &acl,
		Bucket:        &b.nameWithPrefix,
		Body:          obj.data,
		ContentLength: &size,
		ContentType:   pointers.String(obj.FileType()),
		Key:           pointers.String(path),
	}

	_, err := b.service.client.PutObject(req)
	if err != nil {
		b.service.Errorf("error on `PutObject` operation; bucket=%s; error=%s;", b.nameWithPrefix, err.Error())
	}
	return err
}

// GetObjectByte returns bytes of object from given S3 path.
func (b *Bucket) GetObjectByte(path string) ([]byte, error) {
	r, err := b.getObject(path)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	return buf.Bytes(), err
}

// getObject fetches object from target S3 path
func (b *Bucket) getObject(path string) (io.Reader, error) {
	out, err := b.service.client.GetObject(&SDK.GetObjectInput{
		Bucket: &b.name,
		Key:    &path,
	})
	if err != nil {
		b.service.Errorf("error on `GetObject` operation; bucket=%s; error=%s;", b.name, err.Error())
		return nil, err
	}
	return out.Body, nil
}

// GetURL fetches url of target S3 object.
func (b *Bucket) GetURL(path string) string {
	return fmt.Sprintf("%s/%s/%s", b.endpoint, b.nameWithPrefix, path)
}

// GetSecretURL fetches a url of target S3 object w/ ACL permission.
func (b *Bucket) GetSecretURL(path string) (string, error) {
	return b.GetSecretURLWithExpire(path, b.expireSecond)
}

// GetSecretURLWithExpire fetches a url of target S3 object w/ ACL permission (url expires in `expire` value seconds)
// ** this isn't work **
func (b *Bucket) GetSecretURLWithExpire(path string, expire int) (string, error) {
	req, _ := b.service.client.GetObjectRequest(&SDK.GetObjectInput{
		Bucket: pointers.String(b.nameWithPrefix),
		Key:    pointers.String(path),
	})
	return req.Presign(time.Duration(expire) * time.Second)
}

// HeadObject executes HeadObject operation.
func (b *Bucket) HeadObject(path string) (*SDK.HeadObjectOutput, error) {
	return b.service.client.HeadObject(&SDK.HeadObjectInput{
		Bucket: pointers.String(b.nameWithPrefix),
		Key:    pointers.String(path),
	})
}

// IsExists checks if the given path.
func (b *Bucket) IsExists(path string) bool {
	_, err := b.HeadObject(path)
	return err == nil
}

// DeleteObject deletees the object of target path.
func (b *Bucket) DeleteObject(path string) error {
	_, err := b.service.client.DeleteObject(&SDK.DeleteObjectInput{
		Bucket: pointers.String(b.nameWithPrefix),
		Key:    pointers.String(path),
	})
	if err != nil {
		b.service.Errorf("error on `GetObject` operation; bucket=%s; error=%s;", b.nameWithPrefix, err.Error())
	}
	return err
}
