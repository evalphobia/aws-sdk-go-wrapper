package config

import "os"

const (
	envRegion = "AWS_REGION"

	envEndpoint         = "AWS_ENDPOINT"
	envDynamoDBEndpoint = "AWS_DYNAMODB_ENDPOINT"
	envS3Endpoint       = "AWS_S3_ENDPOINT"
	envSNSEndpoint      = "AWS_SNS_ENDPOINT"
	envSQSEndpoint      = "AWS_SQS_ENDPOINT"
)

// EnvRegion get region from env params
func EnvRegion() string {
	return os.Getenv(envRegion)
}

// EnvEndpoint get endpoint from env params
func EnvEndpoint() string {
	return os.Getenv(envEndpoint)
}

// EnvDynamoDBEndpoint get DynamoDB endpoint from env params
func EnvDynamoDBEndpoint() string {
	return os.Getenv(envDynamoDBEndpoint)
}

// EnvS3Endpoint get S3 endpoint from env params
func EnvS3Endpoint() string {
	return os.Getenv(envS3Endpoint)
}

// EnvSNSEndpoint get SNS endpoint from env params
func EnvSNSEndpoint() string {
	return os.Getenv(envSNSEndpoint)
}

// EnvSQSEndpoint get SQS endpoint from env params
func EnvSQSEndpoint() string {
	return os.Getenv(envSQSEndpoint)
}
