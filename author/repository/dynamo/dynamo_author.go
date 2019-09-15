package dynamo

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/jayvib/clean-architecture/apperr"
	"github.com/jayvib/clean-architecture/author"
	"github.com/jayvib/clean-architecture/model"
)

func New(svc dynamodbiface.DynamoDBAPI) author.Repository {
	return &repository{svc: svc}
}

type repository struct {
	svc dynamodbiface.DynamoDBAPI
}

func (r *repository) GetByID(ctx context.Context, id string) (*model.Author, error) {
	tableName := model.GetAuthorTableName()
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}
	output, err := r.svc.GetItemWithContext(ctx, input)
	if err != nil {
		if ae, ok := err.(awserr.Error); ok {
			return nil, apperr.New(ae.Code(), ae.Message(), ae.OrigErr())
		}
		return nil, apperr.New(apperr.InternalError, err.Error(), err)
	}

	var resAuthor model.Author
	err = dynamodbattribute.UnmarshalMap(output.Item, &resAuthor)
	if err != nil {
		return nil, apperr.New(apperr.InternalError, err.Error(), err)
	}

	return &resAuthor, nil
}

func (r *repository) Store(ctx context.Context, u *model.Author) error {
	if u.ID == "" {
		return apperr.New(apperr.EmptyID, "Author to save has no ID", nil)
	}
	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		aerr := apperr.New(apperr.InternalError, "Can't marshal author to dynamo attribute value.", err)
		apperr.AddInfo(aerr, "Author ID", u.ID)
		return aerr
	}
	tableName := u.TableName()
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}
	_, err = r.svc.PutItemWithContext(ctx, input)
	if err != nil {
		if ae, ok := err.(awserr.Error); ok {
			return apperr.New(ae.Code(), ae.Message(), err)
		}
		return apperr.New(apperr.InternalError, err.Error(), err)
	}
	return nil
}
