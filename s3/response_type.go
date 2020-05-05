package s3

import (
	"time"

	SDK "github.com/aws/aws-sdk-go/service/s3"
)

// CopyObjectResponse contains data from CopyObject.
type CopyObjectResponse struct {
	ETag         string
	LastModified time.Time

	CopySourceVersionID     string
	Expiration              string
	RequestCharged          string
	SSECustomerAlgorithm    string
	SSECustomerKeyMD5       string
	SSEKMSEncryptionContext string
	SSEKMSKeyID             string
	ServerSideEncryption    string
	VersionID               string
}

func NewCopyObjectResponse(out *SDK.CopyObjectOutput) CopyObjectResponse {
	r := CopyObjectResponse{}
	if out == nil {
		return r
	}

	if out.CopySourceVersionId != nil {
		r.CopySourceVersionID = *out.CopySourceVersionId
	}
	if out.Expiration != nil {
		r.Expiration = *out.Expiration
	}
	if out.RequestCharged != nil {
		r.RequestCharged = *out.RequestCharged
	}
	if out.SSECustomerAlgorithm != nil {
		r.SSECustomerAlgorithm = *out.SSECustomerAlgorithm
	}
	if out.SSECustomerKeyMD5 != nil {
		r.SSECustomerKeyMD5 = *out.SSECustomerKeyMD5
	}
	if out.SSEKMSEncryptionContext != nil {
		r.SSEKMSEncryptionContext = *out.SSEKMSEncryptionContext
	}
	if out.SSEKMSKeyId != nil {
		r.SSEKMSKeyID = *out.SSEKMSKeyId
	}
	if out.ServerSideEncryption != nil {
		r.ServerSideEncryption = *out.ServerSideEncryption
	}
	if out.VersionId != nil {
		r.VersionID = *out.VersionId
	}

	if out.CopyObjectResult != nil {
		d := out.CopyObjectResult
		if d.ETag != nil {
			r.ETag = *d.ETag
		}
		if d.LastModified != nil {
			r.LastModified = *d.LastModified
		}
	}
	return r
}

// ListObjectsResponse contains data from ListObjectsV2.
type ListObjectsResponse struct {
	CommonPrefixes        []string
	Contents              []Object
	ContinuationToken     string
	Delimiter             string
	EncodingType          string
	IsTruncated           bool
	KeyCount              int64
	MaxKeys               int64
	Name                  string
	NextContinuationToken string
	Prefix                string
	StartAfter            string
}

func NewListObjectsResponse(out *SDK.ListObjectsV2Output) ListObjectsResponse {
	r := ListObjectsResponse{}
	if out == nil {
		return r
	}

	if out.ContinuationToken != nil {
		r.ContinuationToken = *out.ContinuationToken
	}
	if out.Delimiter != nil {
		r.Delimiter = *out.Delimiter
	}
	if out.EncodingType != nil {
		r.EncodingType = *out.EncodingType
	}
	if out.IsTruncated != nil {
		r.IsTruncated = *out.IsTruncated
	}
	if out.KeyCount != nil {
		r.KeyCount = *out.KeyCount
	}
	if out.MaxKeys != nil {
		r.MaxKeys = *out.MaxKeys
	}
	if out.Name != nil {
		r.Name = *out.Name
	}
	if out.NextContinuationToken != nil {
		r.NextContinuationToken = *out.NextContinuationToken
	}
	if out.Prefix != nil {
		r.Prefix = *out.Prefix
	}
	if out.StartAfter != nil {
		r.StartAfter = *out.StartAfter
	}

	if len(out.CommonPrefixes) != 0 {
		list := make([]string, 0, len(out.CommonPrefixes))
		for _, v := range out.CommonPrefixes {
			if v != nil && v.Prefix != nil {
				list = append(list, *v.Prefix)
			}
		}
		r.CommonPrefixes = list
	}

	if len(out.Contents) != 0 {
		list := make([]Object, len(out.Contents))
		for i, v := range out.Contents {
			list[i] = NewObject(v)
		}
		r.Contents = list
	}

	return r
}

type Object struct {
	ETag             string
	Key              string
	LastModified     time.Time
	Size             int64
	StorageClass     string
	OwnerID          string
	OwnerDisplayName string
}

func NewObject(d *SDK.Object) Object {
	o := Object{}
	if d == nil {
		return o
	}

	if d.ETag != nil {
		o.ETag = *d.ETag
	}
	if d.Key != nil {
		o.Key = *d.Key
	}
	if d.Size != nil {
		o.Size = *d.Size
	}
	if d.LastModified != nil {
		o.LastModified = *d.LastModified
	}
	if d.Owner != nil {
		owner := d.Owner
		if owner.ID != nil {
			o.OwnerID = *owner.ID
		}
		if owner.DisplayName != nil {
			o.OwnerDisplayName = *owner.DisplayName
		}
	}
	return o
}
