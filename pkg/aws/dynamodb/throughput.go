package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	sdk "github.com/aws/aws-sdk-go/service/dynamodb"
)

func newProvisionThroughput(r, w int64) *sdk.ProvisionedThroughput {
	return &sdk.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(r),
		WriteCapacityUnits: aws.Int64(w),
	}
}
