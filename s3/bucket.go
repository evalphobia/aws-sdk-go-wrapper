// S3 Bucket setting, Object manipuration

package s3

import (
	// AWS "github.com/awslabs/aws-sdk-go/aws"
	S3 "github.com/awslabs/aws-sdk-go/gen/s3"
	// "github.com/evalphobia/aws-sdk-go-wrapper/log"

	"bytes"
	"errors"
	"io"
	// "time"
)

const (
	defaultExpireSecond = 180
)

// struct for bucket
type Bucket struct {
	name string
	// objects []S3.PutObjectInput

	client *S3.S3
}

// add object to write spool (w/ public read access)
func (b *Bucket) AddObject(obj *S3Object, path string) {
	b.addObject(obj, path, S3.BucketCannedACLPublicRead)
}

// add object to write spool (w/ ACL permission)
func (b *Bucket) AddSecretObject(obj *S3Object, path string) {
	b.addObject(obj, path, S3.BucketCannedACLAuthenticatedRead)
}

// add object to write spool
func (b *Bucket) addObject(obj *S3Object, path, acl string) {
	// size := obj.Size()
	// in := S3.PutObjectInput{
	// 	ACL:           &acl,
	// 	Bucket:        &b.name,
	// 	Body:          obj.data,
	// 	ContentLength: &size,
	// 	ContentType:   AWS.String(obj.FileType()),
	// 	Key:           AWS.String(path),
	// }
	// b.objects = append(b.objects, in)
}

// put object to server
func (b *Bucket) Put() error {
	var err error = nil
	errStr := ""
	// ファイルの保存
	// for _, obj := range b.objects {
	// 	_, e := b.client.PutObject(&obj)
	// 	if e != nil {
	// 		log.Error("[S3] error on `PutObject` operation, bucket="+b.name, e.Error())
	// 		errStr = errStr + "," + e.Error()
	// 	}
	// }
	if errStr != "" {
		err = errors.New(errStr)
	}
	return err
}

// fetch object from target S3 path
func (b *Bucket) getObject(path string) (io.Reader, error) {
	// req := S3.GetObjectInput{
	// 	Bucket: &b.name,
	// 	Key:    &path,
	// }
	// out, err := b.client.GetObject(&req)
	// if err != nil {
	// 	log.Error("[S3] error on `GetObject` operation, bucket="+b.name, err.Error())
	// }
	// return out.Body, err
	return nil, nil
}

// fetch bytes of object from target S3 path
func (b *Bucket) GetObjectByte(path string) ([]byte, error) {
	r, err := b.getObject(path)
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	return buf.Bytes(), err
}

// fetch url of target S3 object
// (this is same as secret at this time, since it is no method for public url)
func (b *Bucket) GetURL(path string) (string, error) {
	return b.GetSecretURLWithExpire(path, defaultExpireSecond)
}

// fetch url of target S3 object w/ ACL permission （url expires in 3min）
func (b *Bucket) GetSecretURL(path string) (string, error) {
	return b.GetSecretURLWithExpire(path, defaultExpireSecond)
}

// fetch url of target S3 object w/ ACL permission （url expires in `expire` value seconds)
func (b *Bucket) GetSecretURLWithExpire(path string, expire uint64) (string, error) {
	// req, _ := b.client.GetObjectRequest(&S3.GetObjectInput{
	// 	Bucket: AWS.String(b.name),
	// 	Key:    AWS.String(path),
	// })
	// return req.Presign(time.Duration(expire) * time.Second)
	return "", nil
}

// delete object of target path
func (b *Bucket) DeleteObject(path string) error {
	// _, err := b.client.DeleteObject(&S3.DeleteObjectInput{
	// 	Bucket: AWS.String(b.name),
	// 	Key:    AWS.String(path),
	// })
	// if err != nil {
	// 	log.Error("[S3] error on `DeleteObject` operation, bucket="+b.name, err.Error())
	// }
	// return err
	return nil
}
