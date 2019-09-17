package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jayvib/app/model"
)

// TODO: use the MarshalMap directly no need to wrap
func marshalToAttribute(u *model.User) (map[string]*dynamodb.AttributeValue, error) {
	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, err
	}
	return av, nil
}

func MarshalUserToAttributeValue(u *model.User) (map[string]*dynamodb.AttributeValue, error) {
	return marshalToAttribute(u)
}

func unmarshalAttributeValueToUser(v map[string]*dynamodb.AttributeValue) (*model.User, error) {
	var u model.User
	if err := dynamodbattribute.UnmarshalMap(v, &u); err != nil {
		return nil, err
	}
	return &u, nil
}
