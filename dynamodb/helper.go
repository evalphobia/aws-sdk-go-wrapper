package dynamodb

import (
	SDK "github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

func newProvisionedThroughput(read, write int64) *SDK.ProvisionedThroughput {
	return &SDK.ProvisionedThroughput{
		ReadCapacityUnits:  pointers.Long64(read),
		WriteCapacityUnits: pointers.Long64(write),
	}
}
