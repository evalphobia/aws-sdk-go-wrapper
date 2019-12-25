package sqs

import (
	"strconv"
)

// AttributesResponse contains attributes from GetQueueAttributes.
type AttributesResponse struct {
	ApproximateNumberOfMessages           int
	ApproximateNumberOfMessagesDelayed    int
	ApproximateNumberOfMessagesNotVisible int
	CreatedTimestamp                      int
	DelaySeconds                          int
	LastModifiedTimestamp                 int
	MaximumMessageSize                    int
	MessageRetentionPeriod                int
	QueueArn                              string
	ReceiveMessageWaitTimeSeconds         int
	RedrivePolicy                         string
	VisibilityTimeout                     int

	// KmsMasterKeyId                        int
	// KmsDataKeyReusePeriodSeconds          int
	// FifoQueue                             bool
	// ContentBasedDeduplication             int
}

func NewAttributesResponse(apiResponse map[string]*string) AttributesResponse {
	a := AttributesResponse{}
	if apiResponse[AttributeApproximateNumberOfMessages] != nil {
		a.ApproximateNumberOfMessages, _ = strconv.Atoi(*apiResponse[AttributeApproximateNumberOfMessages])
	}
	if apiResponse[AttributeApproximateNumberOfMessagesDelayed] != nil {
		a.ApproximateNumberOfMessagesDelayed, _ = strconv.Atoi(*apiResponse[AttributeApproximateNumberOfMessagesDelayed])
	}
	if apiResponse[AttributeApproximateNumberOfMessagesNotVisible] != nil {
		a.ApproximateNumberOfMessagesNotVisible, _ = strconv.Atoi(*apiResponse[AttributeApproximateNumberOfMessagesNotVisible])
	}
	if apiResponse[AttributeCreatedTimestamp] != nil {
		a.CreatedTimestamp, _ = strconv.Atoi(*apiResponse[AttributeCreatedTimestamp])
	}
	if apiResponse[AttributeDelaySeconds] != nil {
		a.DelaySeconds, _ = strconv.Atoi(*apiResponse[AttributeDelaySeconds])
	}
	if apiResponse[AttributeLastModifiedTimestamp] != nil {
		a.LastModifiedTimestamp, _ = strconv.Atoi(*apiResponse[AttributeLastModifiedTimestamp])
	}
	if apiResponse[AttributeMaximumMessageSize] != nil {
		a.MaximumMessageSize, _ = strconv.Atoi(*apiResponse[AttributeMaximumMessageSize])
	}
	if apiResponse[AttributeMessageRetentionPeriod] != nil {
		a.MessageRetentionPeriod, _ = strconv.Atoi(*apiResponse[AttributeMessageRetentionPeriod])
	}
	if apiResponse[AttributeReceiveMessageWaitTimeSeconds] != nil {
		a.ReceiveMessageWaitTimeSeconds, _ = strconv.Atoi(*apiResponse[AttributeReceiveMessageWaitTimeSeconds])
	}
	if apiResponse[AttributeVisibilityTimeout] != nil {
		a.VisibilityTimeout, _ = strconv.Atoi(*apiResponse[AttributeVisibilityTimeout])
	}

	if apiResponse[AttributeQueueArn] != nil {
		a.QueueArn = *apiResponse[AttributeQueueArn]
	}
	if apiResponse[AttributeRedrivePolicy] != nil {
		a.RedrivePolicy = *apiResponse[AttributeRedrivePolicy]
	}
	return a
}
