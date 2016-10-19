package kinesis

import SDK "github.com/aws/aws-sdk-go/service/kinesis"

// RecordResult is struct for result of `GetRecord` operation.
type RecordResult struct {
	ShardID           string
	Items             []*SDK.Record
	Count             int
	NextShardIterator string
	Behind            int64
}
