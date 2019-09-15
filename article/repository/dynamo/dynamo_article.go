package dynamo

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/jayvib/clean-architecture/apperr"
	"github.com/jayvib/clean-architecture/model"
	apptime "github.com/jayvib/clean-architecture/time"
	"time"
)

func New(svc dynamodbiface.DynamoDBAPI) *Repository {
	return &Repository{svc: svc}
}

type Repository struct {
	svc dynamodbiface.DynamoDBAPI
}

func (r *Repository) Fetch(ctx context.Context, cursor string, num int) (ars []*model.Article, nextCursor string, err error) {
	decodedCursor, err := apptime.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", apperr.New("DecodeErr", err.Error(), err)
	}
	input := &dynamodb.ScanInput{
		TableName:        aws.String(model.GetArticleTableName()),
		FilterExpression: aws.String("#ca >= :d"),
		ExpressionAttributeNames: map[string]*string{
			"#ca": aws.String("created_at"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":d": {S: aws.String(decodedCursor.String())},
		},
		Limit: aws.Int64(int64(num)),
	}

	output, err := r.svc.ScanWithContext(ctx, input)
	if err != nil {
		if ae, ok := err.(awserr.Error); ok {
			return nil, "", apperr.New(ae.Code(), ae.Message(), ae.OrigErr())
		}
		return nil, "", apperr.New(apperr.InternalError, err.Error(), err)
	}
	if aws.Int64Value(output.Count) == 0 {
		return nil, "", apperr.New(apperr.NoItemFound,
			"No item found yet from the database", nil)
	}
	res := make([]*model.Article, 0)
	for _, item := range output.Items {
		var ar model.Article
		dynamodbattribute.UnmarshalMap(item, &ar)
		res = append(res, &ar)
	}

	if output.LastEvaluatedKey != nil {
		for _, a := range res {
			// a numb way
			if _, ok := output.LastEvaluatedKey[a.ID]; ok {
				nextCursor = apptime.EncodeCursor(a.CreatedAt)
			}
		}
	}

	return res, nextCursor, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (ar *model.Article, err error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(model.GetArticleTableName()),
		Key: map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{
				S: aws.String(id),
			},
		},
	}
	res, err := r.svc.GetItemWithContext(ctx, input)
	if err != nil {
		if ae, ok := err.(awserr.Error); ok {
			return nil, apperr.New(ae.Code(), ae.Message(), ae.OrigErr())
		}
		return nil, apperr.New(apperr.InternalError, err.Error(), err)
	}

	if res != nil && res.Item == nil {
		err = apperr.New(apperr.NoItemFound, fmt.Sprintf("No item found for ID: %s", id), nil)
		err = apperr.AddInfos(err, "ID", id, "Operation", "GetItemWithContext")
		return nil, err
	}

	var resArticle model.Article
	err = dynamodbattribute.UnmarshalMap(res.Item, &resArticle)
	if err != nil {
		return nil, apperr.New(apperr.InternalError, err.Error(), nil)
	}
	return &resArticle, nil
}

func (r *Repository) GetByTitle(ctx context.Context, title string) (ar *model.Article, err error) {
	tableName := model.GetArticleTableName()
	input := &dynamodb.ScanInput{
		TableName:        aws.String(tableName),
		FilterExpression: aws.String("title= :t"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {
				S: aws.String(title),
			},
		},
	}
	output, err := r.svc.ScanWithContext(ctx, input)
	if err != nil {
		if ae, ok := err.(awserr.Error); ok {
			return nil, apperr.New(ae.Code(), ae.Message(), err)
		}
		aerr := apperr.New(apperr.InternalError, err.Error(), err)
		apperr.AddInfos(aerr, "Title", title)
		return nil, aerr
	}
	if aws.Int64Value(output.Count) < 1 {
		return nil, apperr.New(apperr.NoItemFound, "Item not exist yet in Dynamo", nil)
	}

	var article model.Article
	err = dynamodbattribute.UnmarshalMap(output.Items[0], &article)
	if err != nil {
		return nil, apperr.New(apperr.InternalError, err.Error(), err)
	}

	return &article, nil
}

func (r *Repository) Update(ctx context.Context, ar *model.Article) error {
	tableName := model.GetArticleTableName()
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(ar.ID),
			},
		},

		// How About the nested??
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {
				S: aws.String(ar.Title),
			},
			":c": {
				S: aws.String(ar.Content),
			},
			":ua": {
				S: aws.String(ar.UpdatedAt.Format(time.RFC3339)),
			},
		},
		UpdateExpression: aws.String("set title = :t, content = :c, updated_at = :ua"),
		ReturnValues:     aws.String("NONE"),
	}
	_, err := r.svc.UpdateItemWithContext(ctx, input)
	if err != nil {
		return handleErr(err, apperr.InternalError)
	}
	return nil
}
func (r *Repository) Store(ctx context.Context, ar *model.Article) error {
	av, err := dynamodbattribute.MarshalMap(ar)
	if err != nil {
		return handleErr(err, apperr.InternalError)
	}

	tableName := ar.TableName()
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}

	_, err = r.svc.PutItemWithContext(ctx, input)
	if err != nil {
		handleErr(err, apperr.InternalError)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}
	_, err := r.svc.DeleteItemWithContext(ctx, input)
	if err != nil {
		return handleErr(err, apperr.InternalError)
	}
	return nil
}

func handleErr(err error, errCode string, infos ...string) error {
	if ae, ok := err.(awserr.Error); ok {
		aerr := apperr.New(ae.Code(), ae.Message(), ae.OrigErr())
		if infos != nil {
			apperr.AddInfos(aerr, infos...)
		}
		return aerr
	}
	return apperr.New(errCode, err.Error(), err)
}
