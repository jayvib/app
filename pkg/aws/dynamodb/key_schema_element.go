package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	sdk "github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	KeyTypeHash  = "HASH"
	KeyTypeRange = "RANGE"
)

func NewKeySchemaElement(keyName, keyType string) *sdk.KeySchemaElement {
	return &sdk.KeySchemaElement{
		AttributeName: aws.String(keyName),
		KeyType:       aws.String(keyType),
	}
}

// NewHashKey creates a key schema element for hash key.
func NewHashKey(keyName string) *sdk.KeySchemaElement {
	return NewKeySchemaElement(keyName, KeyTypeHash)
}

// NewRangeKey creates a key schema element for range key.
func NewRangeKey(keyName string) *sdk.KeySchemaElement {
	return NewKeySchemaElement(keyName, KeyTypeRange)
}

// NewKeySchemeElements will control the number of key schema element
// that will be included within the array. So basically there are
// maximum of two items in the array of the key schema which are
// the partition key and the sort key.
func NewKeySchemaElements(elements ...*sdk.KeySchemaElement) []*sdk.KeySchemaElement {
	if len(elements) > 1 {
		keySchemas := make([]*sdk.KeySchemaElement, 2, 2)
		keySchemas[0] = elements[0]
		keySchemas[1] = elements[1]
		return keySchemas
	}
	keySchema := make([]*sdk.KeySchemaElement, 1, 1)
	keySchema[0] = elements[0]
	return keySchema
}
