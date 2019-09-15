package aws

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func RetriableError(err error) bool {
	if awsErr, ok := err.(awserr.Error); ok {
		if awsErr.Code() == dynamodb.ErrCodeProvisionedThroughputExceededException {
			return true
		}
	}
	return false
}
