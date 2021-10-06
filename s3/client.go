package s3

import (
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	SDK "github.com/aws/aws-sdk-go/service/s3"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
	"github.com/evalphobia/aws-sdk-go-wrapper/private/errors"
	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

const (
	serviceName = "s3"
)

// S3 has S3 client and bucket list.
type S3 struct {
	client   *SDK.S3
	endpoint string

	logger log.Logger
	prefix string

	bucketsMu sync.RWMutex
	buckets   map[string]*Bucket
}

// New returns initialized *S3.
func New(conf config.Config) (*S3, error) {
	sess, err := conf.Session()
	if err != nil {
		return nil, err
	}

	svc := NewFromSession(sess)
	svc.prefix = conf.DefaultPrefix
	return svc, nil
}

// NewFromSession returns initialized *S3 from aws.Session.
func NewFromSession(sess *session.Session) *S3 {
	cli := SDK.New(sess)
	return &S3{
		client:   cli,
		endpoint: cli.ClientInfo.Endpoint,
		logger:   log.DefaultLogger,
		buckets:  make(map[string]*Bucket),
	}
}

// SetLogger sets logger.
func (svc *S3) SetLogger(logger log.Logger) {
	svc.logger = logger
}

// GetBucket gets S3 bucket.
func (svc *S3) GetBucket(bucket string) (*Bucket, error) {
	bucketName := svc.prefix + bucket

	// get the bucket from cache
	svc.bucketsMu.RLock()
	b, ok := svc.buckets[bucketName]
	svc.bucketsMu.RUnlock()
	if ok {
		return b, nil
	}

	// get the bucket from AWS api.
	_, err := svc.client.GetBucketLocation(&SDK.GetBucketLocationInput{
		Bucket: pointers.String(bucketName),
	})
	if err != nil {
		svc.Errorf("error on `GetQueueURL` operation; bueckt=%s; error=%s;", bucketName, err.Error())
		return nil, err
	}

	b = NewBucket(svc, bucket)
	svc.bucketsMu.Lock()
	svc.buckets[bucketName] = b
	svc.bucketsMu.Unlock()
	return b, nil
}

// IsExistBucket checks if the Bucket already exists or not.
func (svc *S3) IsExistBucket(name string) (bool, error) {
	bucketName := svc.prefix + name
	// get the bucket from AWS api.
	data, err := svc.client.GetBucketLocation(&SDK.GetBucketLocationInput{
		Bucket: pointers.String(bucketName),
	})

	switch {
	case isNonSuchBucketError(err):
		return false, nil
	case err != nil:
		svc.Errorf("error on `GetQueueUrl` operation; queue=%s; error=%s", name, err.Error())
		return false, err
	case data != nil:
		return true, nil
	default:
		return false, nil
	}
}

func isNonSuchBucketError(err error) bool {
	const errNonSuchBucket = "NoSuchBucket: "
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), errNonSuchBucket)
}

// CreateBucket creates new S3 bucket.
func (svc *S3) CreateBucket(in *SDK.CreateBucketInput) error {
	data, err := svc.client.CreateBucket(in)
	if err != nil {
		svc.Errorf("error on `CreateBucket` operation; bucket=%s; error=%s;", *in.Bucket, err.Error())
		return err
	}

	svc.Infof("success on `CreateBucket` operation; bucket=%s; data=%s;", *in.Bucket, data.String())
	return nil
}

// CreateBucketWithName creates new S3 bucket by given name.
func (svc *S3) CreateBucketWithName(name string) error {
	bucketName := svc.prefix + name
	return svc.CreateBucket(&SDK.CreateBucketInput{
		Bucket: pointers.String(bucketName),
	})
}

// ForceDeleteBucket deletes S3 bucket by given name.
func (svc *S3) ForceDeleteBucket(name string) error {
	bucketName := svc.prefix + name
	_, err := svc.client.DeleteBucket(&SDK.DeleteBucketInput{
		Bucket: pointers.String(bucketName),
	})
	if err != nil {
		svc.Errorf("error on `DeleteBucket` operation; bucket=%s; error=%s;", name, err.Error())
		return err
	}

	svc.Infof("success on `DeleteBucket` operation; bucket=%s;", name)
	return nil
}

// CopyObject executes `CopyObject` operation.
func (svc *S3) CopyObject(req CopyObjectRequest) (CopyObjectResponse, error) {
	out, err := svc.copyObject(req.ToInput())
	return NewCopyObjectResponse(out), err
}

func (svc *S3) copyObject(input *SDK.CopyObjectInput) (*SDK.CopyObjectOutput, error) {
	return svc.client.CopyObject(input)
}

// Infof logging information.
func (svc *S3) Infof(format string, v ...interface{}) {
	svc.logger.Infof(serviceName, format, v...)
}

// Errorf logging error information.
func (svc *S3) Errorf(format string, v ...interface{}) {
	svc.logger.Errorf(serviceName, format, v...)
}

func newErrors() *errors.Errors {
	return errors.NewErrors(serviceName)
}
