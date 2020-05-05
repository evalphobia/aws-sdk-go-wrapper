package s3

import (
	"fmt"
	"time"

	SDK "github.com/aws/aws-sdk-go/service/s3"
	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

// CopyObjectRequest has parameters for `CopyObject` operation.
type CopyObjectRequest struct {
	SrcBucket  string
	SrcPath    string
	DestBucket string
	DestPath   string
	// if true, add prefix to bucket name
	UseSamePrefix bool

	// optional params
	// ref: https://docs.aws.amazon.com/AmazonS3/latest/API/API_CopyObject.html
	ACL                            string
	CacheControl                   string
	ContentDisposition             string
	ContentEncoding                string
	ContentLanguage                string
	ContentType                    string
	CopySourceIfMatch              string
	CopySourceIfModifiedSince      time.Time
	CopySourceIfNoneMatch          string
	CopySourceIfUnmodifiedSince    time.Time
	CopySourceSSECustomerAlgorithm string
	CopySourceSSECustomerKey       string
	CopySourceSSECustomerKeyMD5    string
	Expires                        time.Time
	GrantFullControl               string
	GrantRead                      string
	GrantReadACP                   string
	GrantWriteACP                  string
	Metadata                       map[string]string
	MetadataDirective              string
	ObjectLockLegalHoldStatus      string
	ObjectLockMode                 string
	ObjectLockRetainUntilDate      time.Time
	RequestPayer                   string
	SSECustomerAlgorithm           string
	SSECustomerKey                 string
	SSECustomerKeyMD5              string
	SSEKMSEncryptionContext        string
	SSEKMSKeyID                    string
	ServerSideEncryption           string
	StorageClass                   string
	Tagging                        string
	TaggingDirective               string
	WebsiteRedirectLocation        string
}

func (r CopyObjectRequest) ToInput() *SDK.CopyObjectInput {
	in := &SDK.CopyObjectInput{}

	if r.SrcBucket != "" && r.SrcPath != "" {
		in.SetCopySource(fmt.Sprintf("/%s/%s", r.SrcBucket, r.SrcPath))
	}
	if r.DestBucket != "" {
		in.SetBucket(r.DestBucket)
	}
	if r.DestPath != "" {
		in.SetKey(r.DestPath)
	}

	if r.ACL != "" {
		in.SetACL(r.ACL)
	}
	if r.CacheControl != "" {
		in.SetCacheControl(r.CacheControl)
	}
	if r.ContentDisposition != "" {
		in.SetContentDisposition(r.ContentDisposition)
	}
	if r.ContentEncoding != "" {
		in.SetContentEncoding(r.ContentEncoding)
	}
	if r.ContentLanguage != "" {
		in.SetContentLanguage(r.ContentLanguage)
	}
	if r.ContentType != "" {
		in.SetContentType(r.ContentType)
	}
	if r.CopySourceIfMatch != "" {
		in.SetCopySourceIfMatch(r.CopySourceIfMatch)
	}
	if r.CopySourceIfNoneMatch != "" {
		in.SetCopySourceIfNoneMatch(r.CopySourceIfNoneMatch)
	}
	if r.CopySourceSSECustomerAlgorithm != "" {
		in.SetCopySourceSSECustomerAlgorithm(r.CopySourceSSECustomerAlgorithm)
	}
	if r.CopySourceSSECustomerKey != "" {
		in.SetCopySourceSSECustomerKey(r.CopySourceSSECustomerKey)
	}
	if r.CopySourceSSECustomerKeyMD5 != "" {
		in.SetCopySourceSSECustomerKeyMD5(r.CopySourceSSECustomerKeyMD5)
	}
	if r.GrantFullControl != "" {
		in.SetGrantFullControl(r.GrantFullControl)
	}
	if r.GrantRead != "" {
		in.SetGrantRead(r.GrantRead)
	}
	if r.GrantReadACP != "" {
		in.SetGrantReadACP(r.GrantReadACP)
	}
	if r.GrantWriteACP != "" {
		in.SetGrantWriteACP(r.GrantWriteACP)
	}
	if r.MetadataDirective != "" {
		in.SetMetadataDirective(r.MetadataDirective)
	}
	if r.ObjectLockLegalHoldStatus != "" {
		in.SetObjectLockLegalHoldStatus(r.ObjectLockLegalHoldStatus)
	}
	if r.ObjectLockMode != "" {
		in.SetObjectLockMode(r.ObjectLockMode)
	}
	if r.RequestPayer != "" {
		in.SetRequestPayer(r.RequestPayer)
	}
	if r.SSECustomerAlgorithm != "" {
		in.SetSSECustomerAlgorithm(r.SSECustomerAlgorithm)
	}
	if r.SSECustomerKey != "" {
		in.SetSSECustomerKey(r.SSECustomerKey)
	}
	if r.SSECustomerKeyMD5 != "" {
		in.SetSSECustomerKeyMD5(r.SSECustomerKeyMD5)
	}
	if r.SSEKMSEncryptionContext != "" {
		in.SetSSEKMSEncryptionContext(r.SSEKMSEncryptionContext)
	}
	if r.SSEKMSKeyID != "" {
		in.SetSSEKMSKeyId(r.SSEKMSKeyID)
	}
	if r.ServerSideEncryption != "" {
		in.SetServerSideEncryption(r.ServerSideEncryption)
	}
	if r.StorageClass != "" {
		in.SetStorageClass(r.StorageClass)
	}
	if r.Tagging != "" {
		in.SetTagging(r.Tagging)
	}
	if r.TaggingDirective != "" {
		in.SetTaggingDirective(r.TaggingDirective)
	}
	if r.WebsiteRedirectLocation != "" {
		in.SetWebsiteRedirectLocation(r.WebsiteRedirectLocation)
	}

	if !r.CopySourceIfModifiedSince.IsZero() {
		in.SetCopySourceIfModifiedSince(r.CopySourceIfModifiedSince)
	}
	if !r.CopySourceIfUnmodifiedSince.IsZero() {
		in.SetCopySourceIfUnmodifiedSince(r.CopySourceIfUnmodifiedSince)
	}
	if !r.Expires.IsZero() {
		in.SetExpires(r.Expires)
	}
	if !r.ObjectLockRetainUntilDate.IsZero() {
		in.SetObjectLockRetainUntilDate(r.ObjectLockRetainUntilDate)
	}

	if len(r.Metadata) != 0 {
		m := make(map[string]*string, len(r.Metadata))
		for k, v := range r.Metadata {
			m[k] = pointers.String(v)
		}
		in.SetMetadata(m)
	}
	return in
}

// ListObjectsRequest has parameters for `ListObjectsV2` operation.
type ListObjectsRequest struct {
	Bucket string

	// optional
	ContinuationToken string
	Delimiter         string
	EncodingType      string
	FetchOwner        bool
	MaxKeys           int64
	Prefix            string
	RequestPayer      string
	StartAfter        string
}

func (r ListObjectsRequest) ToInput() *SDK.ListObjectsV2Input {
	in := &SDK.ListObjectsV2Input{}
	if r.Bucket != "" {
		in.SetBucket(r.Bucket)
	}

	if r.ContinuationToken != "" {
		in.SetContinuationToken(r.ContinuationToken)
	}
	if r.Delimiter != "" {
		in.SetDelimiter(r.Delimiter)
	}
	if r.EncodingType != "" {
		in.SetEncodingType(r.EncodingType)
	}
	if r.FetchOwner {
		in.SetFetchOwner(r.FetchOwner)
	}
	if r.MaxKeys != 0 {
		in.SetMaxKeys(r.MaxKeys)
	}
	if r.Prefix != "" {
		in.SetPrefix(r.Prefix)
	}
	if r.RequestPayer != "" {
		in.SetRequestPayer(r.RequestPayer)
	}
	if r.StartAfter != "" {
		in.SetStartAfter(r.StartAfter)
	}
	return in
}
