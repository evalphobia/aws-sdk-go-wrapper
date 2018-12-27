package kms

import (
	SDK "github.com/aws/aws-sdk-go/service/kms"

	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

// Tag is key-value data struct.
type Tag struct {
	Key   string
	Value string
}

func (t Tag) Tag() *SDK.Tag {
	return &SDK.Tag{
		TagKey:   pointers.String(t.Key),
		TagValue: pointers.String(t.Value),
	}
}
