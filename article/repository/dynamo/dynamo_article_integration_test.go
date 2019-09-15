// +build integration,dynamo

package dynamo_test

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jayvib/clean-architecture/article"
	"github.com/jayvib/clean-architecture/article/repository/dynamo"
	"github.com/jayvib/clean-architecture/model"
	apptime "github.com/jayvib/clean-architecture/time"
	"github.com/jayvib/clean-architecture/utils/generateutil"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
	"time"
)

const (
	// Explicit setup to avoid any issue in the production
	defaultEndpoint  = "http://localhost:8000"
	defaultRegion    = "us-east-1"
	defaultAccessKey = "access"
	defaultSecretKey = "secret"
)

var repo article.Repository

var articles = []*model.Article{
	{
		ID:      generateutil.GenerateID("article"),
		Title:   "Pirate King",
		Content: "Luffy will be the next Pirate King!",
		Author: &model.Author{
			ID:   "unqieuid",
			Name: "Luffy Monkey",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:      generateutil.GenerateID("article"),
		Title:   "Zoro The Great Swordsman",
		Content: "Zoro will be the next Swordsman Master",
		Author: &model.Author{
			ID:   "uniquid",
			Name: "Roronoa Zoro",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

func TestMain(m *testing.M) {
	var tdown func()
	var err error
	repo, tdown, err = setup()
	if err != nil {
		log.Fatal(err)
	}
	exitVal := m.Run()
	tdown()
	os.Exit(exitVal)
}

func TestIntegration_GetByID(t *testing.T) {
	articleRes, err := repo.GetByID(context.Background(), articles[0].ID)
	require.NoError(t, err)
	assert.Equal(t, articles[0].ID, articleRes.ID)
	assert.Equal(t, articles[0].Title, articleRes.Title)
	assert.Equal(t, articles[0].Content, articleRes.Content)
	//t.Logf("%#v", articleRes)
}

func TestIntegration_Delete(t *testing.T) {

}

func TestIntegration_GetByTitle(t *testing.T) {
	articleRes, err := repo.GetByTitle(context.Background(), articles[0].Title)
	require.NoError(t, err)
	assert.Equal(t, articles[0].ID, articleRes.ID)
	assert.Equal(t, articles[0].Title, articleRes.Title)
	assert.Equal(t, articles[0].Content, articleRes.Content)
}

func TestIntegration_Fetch(t *testing.T) {
	cursor := apptime.EncodeCursor(time.Now().Add(-1 * time.Hour))
	res, nextCursor, err := repo.Fetch(context.Background(), cursor, 2)
	require.NoError(t, err)
	assert.Empty(t, nextCursor)
	assert.Len(t, res, 2)
	//t.Logf("%#v", res[0])
}

func TestIntegration_Store(t *testing.T) {
	newArticle := articles[0]
	newArticle.ID = generateutil.GenerateID("article")
	newArticle.Title = "New title"
	err := repo.Store(context.Background(), newArticle)
	require.NoError(t, err)

	resArticle, err := repo.GetByID(context.Background(), newArticle.ID)
	require.NoError(t, err)
	assert.Equal(t, newArticle.Title, resArticle.Title)
	assert.Equal(t, newArticle.Content, resArticle.Content)
}

func TestIntegration_Update(t *testing.T) {
	newArticle := articles[0]
	newArticle.Title = "Gear 4"
	err := repo.Update(context.Background(), newArticle)
	require.NoError(t, err)

	resArticle, err := repo.GetByID(context.Background(), newArticle.ID)
	require.NoError(t, err)
	assert.Equal(t, newArticle.Title, resArticle.Title)
	assert.Equal(t, newArticle.Content, resArticle.Content)
}

func setup() (repo article.Repository, tearDown func(), err error) {
	svc, err := newService()
	if err != nil {
		return nil, nil, err
	}
	logrus.Println("setup: creating table")
	err = createTable(svc)
	if err != nil {
		return nil, nil, err
	}
	logrus.Println("setup: adding items")
	err = addItemsToTable(svc)
	return dynamo.New(svc), func() {
		logrus.Println("teardown: cleaning up resources.")
		err = teardown(svc)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Printf("Tear down error")
		}
	}, nil
}

func newService() (*dynamodb.DynamoDB, error) {
	cred := credentials.NewStaticCredentials(defaultAccessKey, defaultSecretKey, "")
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(defaultRegion),
			Endpoint:    aws.String(defaultEndpoint),
			Credentials: cred,
		},
	})
	if err != nil {
		return nil, err
	}
	svc := dynamodb.New(sess)
	return svc, nil
}

func createTable(svc *dynamodb.DynamoDB) error {
	_, err := svc.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String("article"),
	})
	return err
}

func getItem(svc *dynamodb.DynamoDB, id string) (*model.Article, error) {
	tableName := model.GetArticleTableName()
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}

	output, err := svc.GetItem(input)
	if err != nil {
		return nil, err
	}

	if output.Item == nil {
		return nil, errors.New("not found")
	}

	var articleRes model.Article
	err = dynamodbattribute.UnmarshalMap(output.Item, &articleRes)
	if err != nil {
		return nil, err
	}
	return &articleRes, nil
}

func addItemsToTable(svc *dynamodb.DynamoDB) error {
	for _, article := range articles {
		av, err := dynamodbattribute.MarshalMap(article)
		if err != nil {
			return err
		}
		input := &dynamodb.PutItemInput{
			TableName: aws.String("article"),
			Item:      av,
		}

		_, err = svc.PutItem(input)
		if err != nil {
			return err
		}
	}
	return nil
}

func teardown(svc *dynamodb.DynamoDB) error {
	_, err := svc.DeleteTable(&dynamodb.DeleteTableInput{
		TableName: aws.String("article"),
	})
	return err
}
