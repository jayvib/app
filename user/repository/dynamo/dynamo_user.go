package dynamo

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	apperror "github.com/jayvib/clean-architecture/apperr"
	"github.com/jayvib/clean-architecture/model"
	"github.com/jayvib/clean-architecture/user"
	"github.com/sirupsen/logrus"
	"time"
)

func New(db dynamodbiface.DynamoDBAPI) user.Repository {
	return &UserRepository{db: db}
}

type UserRepository struct {
	db dynamodbiface.DynamoDBAPI
}

func (u *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	tablename := model.GetUserTableName()
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tablename),
		// TODO: Would it be nice if ma automate nako ang
		// pag convert to ma[string]*dynamodb.AttributeValue?
		Key: map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{
				S: aws.String(id),
			},
		},
	}

	res, err := u.db.GetItemWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, apperror.ItemNotFound
	}
	//logrus.Printf("%#v", res)
	resUser, err := unmarshalAttributeValueToUser(res.Item)
	if err != nil {
		return nil, err
	}
	return resUser, nil
}

func (u *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	tableName := model.GetUserTableName()
	input := &dynamodb.ScanInput{
		TableName:        aws.String(tableName),
		FilterExpression: aws.String("email= :e"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":e": &dynamodb.AttributeValue{
				S: aws.String(email),
			},
		},
	}
	output, err := u.db.ScanWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	if len(output.Items) == 0 {
		return nil, apperror.ItemNotFound
	}
	usr, err := unmarshalAttributeValueToUser(output.Items[0])
	return usr, err
}

func (u *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	tableName := model.GetUserTableName()
	input := &dynamodb.ScanInput{
		TableName:        aws.String(tableName),
		FilterExpression: aws.String("username= :un"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":un": &dynamodb.AttributeValue{
				S: aws.String(username),
			},
		},
	}
	output, err := u.db.ScanWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	if len(output.Items) == 0 {
		return nil, apperror.ItemNotFound
	}
	return unmarshalAttributeValueToUser(output.Items[0])
}

func (u *UserRepository) Update(ctx context.Context, user *model.User) error {
	if user.ID == "" {
		return apperror.EmptyItemID
	}

	tablename := user.TableName()

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tablename),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(user.ID),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":fn": {
				S: aws.String(user.Firstname),
			},
			":ln": {
				S: aws.String(user.Lastname),
			},
			":e": {
				S: aws.String(user.Email),
			},
			":p": {
				S: aws.String(user.Password),
			},
			":ua": {
				S: aws.String(user.UpdatedAt.Format(time.RFC3339)),
			},
		},
		UpdateExpression: aws.String("SET firstname = :fn, lastname = :ln, email = :e, password = :p, updated_at = :ua"),
	}

	_, err := u.db.UpdateItemWithContext(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) Store(ctx context.Context, user *model.User) error {
	if user.ID == "" {
		return apperror.EmptyItemID
	}
	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"av": av,
	}).Debug()

	tablename := user.TableName()
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tablename),
		Item:      av,
	}

	if ctx == nil {
		ctx = context.Background()
	}
	_, err = u.db.PutItemWithContext(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) Delete(ctx context.Context, id string) error {

	if ctx == nil {
		ctx = context.Background()
	}

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}

	_, err := u.db.DeleteItemWithContext(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
