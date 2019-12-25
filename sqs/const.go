package sqs

// Attribute names for SQS.
// ref: https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_GetQueueAttributes.html
const (
	AttributeAll                                   = "All"
	AttributeApproximateNumberOfMessages           = "ApproximateNumberOfMessages"
	AttributeApproximateNumberOfMessagesDelayed    = "ApproximateNumberOfMessagesDelayed"
	AttributeApproximateNumberOfMessagesNotVisible = "ApproximateNumberOfMessagesNotVisible"
	AttributeCreatedTimestamp                      = "CreatedTimestamp"
	AttributeDelaySeconds                          = "DelaySeconds"
	AttributeLastModifiedTimestamp                 = "LastModifiedTimestamp"
	AttributeMaximumMessageSize                    = "MaximumMessageSize"
	AttributeMessageRetentionPeriod                = "MessageRetentionPeriod"
	AttributePolicy                                = "Policy"
	AttributeQueueArn                              = "QueueArn"
	AttributeReceiveMessageWaitTimeSeconds         = "ReceiveMessageWaitTimeSeconds"
	AttributeRedrivePolicy                         = "RedrivePolicy"
	AttributeVisibilityTimeout                     = "VisibilityTimeout"
	AttributeKmsMasterKeyId                        = "KmsMasterKeyId"
	AttributeKmsDataKeyReusePeriodSeconds          = "KmsDataKeyReusePeriodSeconds"
	AttributeFifoQueue                             = "FifoQueue"
	AttributeContentBasedDeduplication             = "ContentBasedDeduplication"

	AttributeRedrivePolicyDeadLetterTargetArn = "deadLetterTargetArn"
	AttributeRedrivePolicyMaxReceiveCount     = "maxReceiveCount"
)
