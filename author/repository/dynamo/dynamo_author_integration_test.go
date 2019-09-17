// +build integration,dynamo

package dynamo_test

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jayvib/app/author"
	"github.com/jayvib/app/author/repository/dynamo"
	"github.com/jayvib/app/model"
	"github.com/jayvib/app/utils/generateutil"
	"github.com/magiconair/properties/assert"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
	"time"
)

const (
	defaultEndpoint  = "http://localhost:8000"
	defaultRegion    = "us-east-1"
	defaultAccessKey = "access"
	defaultSecretKey = "secret"
)

var repo author.Repository

var mockAuthor = &model.Author{
	ID:        generateutil.GenerateID("author"),
	Name:      "Luffy Monkey",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
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
	t.Run("Success", func(t *testing.T) {
		au, err := repo.GetByID(context.Background(), mockAuthor.ID)
		require.NoError(t, err)
		assert.Equal(t, mockAuthor.Name, au.Name)
	})
}

func TestIntegration_Store(t *testing.T) {
	copyAuthor := *mockAuthor
	copyAuthor.ID = generateutil.GenerateID("author")
	err := repo.Store(context.Background(), &copyAuthor)
	require.NoError(t, err)

	// Get the author
	res, err := repo.GetByID(context.Background(), copyAuthor.ID)
	require.NoError(t, err)
	assert.Equal(t, copyAuthor.Name, res.Name)
}

func setup() (repo author.Repository, tDown func(), err error) {
	svc, err := newDb()
	if err != nil {
		return nil, nil, err
	}
	err = createTable(svc)
	if err != nil {
		return nil, nil, err
	}
	err = putItem(svc)
	if err != nil {
		return nil, nil, err
	}
	return dynamo.New(svc), func() {
		err = tearDown(svc)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Printf("Tear down error")
		}
	}, nil
}

func newDb() (*dynamodb.DynamoDB, error) {
	cred := credentials.NewStaticCredentials(defaultAccessKey, defaultSecretKey, "")
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String("us-east-1"),
			Endpoint:    aws.String("http://localhost:8000"),
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
		TableName: aws.String("author"),
	})
	return err
}

func putItem(svc *dynamodb.DynamoDB) error {
	av, err := dynamodbattribute.MarshalMap(mockAuthor)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String("author"),
		Item:      av,
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}
	return nil
}

func tearDown(svc *dynamodb.DynamoDB) error {
	_, err := svc.DeleteTable(&dynamodb.DeleteTableInput{
		TableName: aws.String("author"),
	})
	return err
}
