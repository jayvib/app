package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	sdk "github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	ProjectionTypeAll      = "ALL"
	ProjectionTypeKeysOnly = "KEYS_ONLY"
)

func NewLocalSecondaryIndex(indexName string, schema []*sdk.KeySchemaElement, projection ...string) *sdk.LocalSecondaryIndex {
	var proj string
	switch {
	case len(projection) == 1:
		proj = projection[0]
	default:
		proj = ProjectionTypeAll
	}
	return &sdk.LocalSecondaryIndex{
		IndexName:  aws.String(indexName),
		KeySchema:  schema,
		Projection: &sdk.Projection{ProjectionType: aws.String(proj)},
	}
}

func NewGlobalSecondaryIndex(indexName string, schema []*sdk.KeySchemaElement, provision *sdk.ProvisionedThroughput, projection ...string) *sdk.GlobalSecondaryIndex {
	var proj string
	switch {
	case len(projection) == 1:
		proj = projection[0]
	default:
		proj = ProjectionTypeAll
	}
	return &sdk.GlobalSecondaryIndex{
		IndexName:             aws.String(indexName),
		KeySchema:             schema,
		Projection:            &sdk.Projection{ProjectionType: aws.String(proj)},
		ProvisionedThroughput: provision,
	}
}
