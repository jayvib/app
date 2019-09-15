package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	sdk "github.com/aws/aws-sdk-go/service/dynamodb"
)

func NewNumberAttributeDefinition(attrName string) *sdk.AttributeDefinition {
	return newAttributeDefinition(attrName, "N")
}

func NewStringAttributeDefinition(attrName string) *sdk.AttributeDefinition {
	return newAttributeDefinition(attrName, "S")
}

func newAttributeDefinition(attrName, attrType string) *sdk.AttributeDefinition {
	return &sdk.AttributeDefinition{
		AttributeName: aws.String(attrName),
		AttributeType: aws.String(attrType),
	}
}
