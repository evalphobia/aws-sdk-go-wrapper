package kinesis

import (
	"fmt"

	SDK "github.com/aws/aws-sdk-go/service/kinesis"

	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

// Stream is Kinesis Stream wrapper struct.
type Stream struct {
	service *Kinesis

	name           string
	nameWithPrefix string
	shardIDs       []string
}

// NewStream returns initialized *Stream.
func NewStream(svc *Kinesis, name string) (*Stream, error) {
	streamName := svc.prefix + name
	s := &Stream{
		service:        svc,
		name:           name,
		nameWithPrefix: streamName,
	}

	// get the stream description from AWS api.
	desc, err := svc.client.DescribeStream(&SDK.DescribeStreamInput{
		StreamName: pointers.String(streamName),
	})
	switch {
	case err != nil:
		svc.Errorf("error on `DescribeStream` operation; stream=%s; error=%s;", streamName, err.Error())
		return nil, err
	case desc == nil, desc.StreamDescription == nil:
		err := fmt.Errorf("error on `DescribeStream`, empty response ; stream=%s;", streamName)
		svc.Errorf(err.Error())
		return nil, err
	}

	// extract shrad id list
	ids := make([]string, len(desc.StreamDescription.Shards))
	for i, shard := range desc.StreamDescription.Shards {
		ids[i] = *(shard.ShardId)
	}
	s.shardIDs = ids
	return s, nil
}

// GetShardIDs returns shard id list of the stream.
func (s *Stream) GetShardIDs() (shardIDs []string, err error) {
	if len(s.shardIDs) != 0 {
		return s.shardIDs, nil
	}

	res, err := s.service.client.DescribeStream(&SDK.DescribeStreamInput{
		StreamName: pointers.String(s.nameWithPrefix),
	})
	switch {
	case err != nil:
		s.service.Errorf("error on `DescribeStream` operation; stream=%s; error=%s;", s.nameWithPrefix, err.Error())
		return nil, err
	case res == nil, res.StreamDescription == nil:
		return nil, fmt.Errorf("cannot find StreamDescription; stream=%s;", s.nameWithPrefix)
	}

	shardIDs = make([]string, len(res.StreamDescription.Shards))
	for i, shard := range res.StreamDescription.Shards {
		shardIDs[i] = *shard.ShardId
	}
	return shardIDs, nil
}

// GetLatestRecords gets records from all of the shards.
func (s *Stream) GetLatestRecords() ([]RecordResult, error) {
	shardIDs, err := s.GetShardIDs()
	if err != nil {
		return nil, err
	}

	var list []RecordResult
	for _, sid := range shardIDs {
		result, err := s.GetRecords(GetCondition{
			ShardID:           sid,
			ShardIteratorType: IteratorTypeLatest,
		})
		if err != nil {
			continue
		}
		list = append(list, result)
	}
	return list, nil
}

// GetRecords gets record fron given condition.
func (s *Stream) GetRecords(cond GetCondition) (RecordResult, error) {
	if cond.ShardIterator == "" {
		shardIter, err := s.getShardIterator(cond.ShardID, cond.ShardIteratorType)
		if err != nil {
			return RecordResult{}, err
		}
		cond.ShardIterator = shardIter
	}

	in := &SDK.GetRecordsInput{
		ShardIterator: pointers.String(cond.ShardIterator),
	}
	if cond.Limit != 0 {
		in.Limit = pointers.Long64(cond.Limit)
	}

	resp, err := s.service.client.GetRecords(in)
	if err != nil {
		s.service.Errorf("error on `GetRecords` operation; stream=%s; error=%s;", s.nameWithPrefix, err.Error())
		return RecordResult{}, err
	}
	return RecordResult{
		ShardID:           cond.ShardID,
		Items:             resp.Records,
		Count:             len(resp.Records),
		Behind:            *resp.MillisBehindLatest,
		NextShardIterator: *resp.NextShardIterator,
	}, nil
}

func (s *Stream) getShardIterator(shardID string, iteratorType IteratorType) (shardIterator string, err error) {
	resp, err := s.service.client.GetShardIterator(&SDK.GetShardIteratorInput{
		ShardId:           pointers.String(shardID),
		ShardIteratorType: pointers.String(iteratorType.String()),
		StreamName:        pointers.String(s.nameWithPrefix),
	})
	if err != nil {
		s.service.Errorf("error on `GetShardIterator` operation; stream=%s; error=%s;", s.nameWithPrefix, err.Error())
		return "", err
	}
	return *resp.ShardIterator, nil
}

// PutRecord puts the given data into stream record.
func (s *Stream) PutRecord(data []byte) error {
	_, err := s.service.client.PutRecord(&SDK.PutRecordInput{
		StreamName:   pointers.String(s.nameWithPrefix),
		PartitionKey: pointers.String(string(data)),
		Data:         data,
	})
	if err != nil {
		s.service.Errorf("error on `PutRecord` operation; stream=%s; error=%s;", s.nameWithPrefix, err.Error())
	}
	return err
}
