// +build integration,dynamo

package dynamo_test

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jayvib/clean-architecture/model"
	"github.com/jayvib/clean-architecture/user"
	"github.com/jayvib/clean-architecture/user/repository/dynamo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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

var repo user.Repository

var expectedUser = &model.User{
	ID:        "uniqueid",
	Firstname: "Luffy",
	Lastname:  "Monkey",
	Email:     "luffy.monkey@onepiece.com",
	Username:  "luffy.monkey",
}

func setup() (repo user.Repository, tdown func(), err error) {
	svc, err := newDb()
	if err != nil {
		return nil, nil, err
	}
	err = createTable(svc)
	if err != nil {
		return nil, nil, err
	}
	return dynamo.New(svc), func() {
		err = teardown(svc)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Printf("Tear down error")
		}
	}, nil
}

func TestMain(m *testing.M) {
	var tdown func()
	var err error
	repo, tdown, err = setup()
	if err != nil {
		log.Fatal(err)
	}
	exitval := m.Run()
	tdown()
	os.Exit(exitval)
}

func TestIntegration_Store(t *testing.T) {
	//logrus.SetLevel(logrus.DebugLevel)
	newUser := *expectedUser
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()
	//t.Logf("%#v", newUser)
	err := repo.Store(context.Background(), &newUser)
	require.NoError(t, err)
	assert.NotEqual(t, expectedUser, newUser)

	// Get the Item
	res, err := repo.GetByID(context.Background(), newUser.ID)
	require.NoError(t, err)

	assert.Equal(t, res.ID, newUser.ID)
	assert.Equal(t, res.Firstname, newUser.Firstname)
}

func TestIntegration_GetByID(t *testing.T) {
	u, err := repo.GetByID(context.Background(), expectedUser.ID)
	require.NoError(t, err)
	assert.Equal(t, expectedUser.Firstname, u.Firstname)
	//t.Logf("%#v", u)
}

func TestIntegration_GetByEmail(t *testing.T) {
	u, err := repo.GetByEmail(context.Background(), expectedUser.Email)
	require.NoError(t, err)
	assert.Equal(t, expectedUser.Email, u.Email)
}

func TestIntegration_GetByUsername(t *testing.T) {
	u, err := repo.GetByUsername(context.Background(), expectedUser.Username)
	require.NoError(t, err)
	assert.Equal(t, expectedUser.Username, u.Username)
}

func TestIntegration_Update(t *testing.T) {
	newUser := expectedUser
	newUser.Email = "luffy.pirateking@onepiece.com"
	newUser.Password = "mysecretpassword"
	err := repo.Update(context.Background(), newUser)
	require.NoError(t, err)

	userRes, err := repo.GetByID(context.Background(), newUser.ID)
	require.NoError(t, err)
	assert.Equal(t, expectedUser.Email, userRes.Email)
	assert.Equal(t, expectedUser.Password, userRes.Password)
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
		TableName: aws.String("user"),
	})
	return err
}

func teardown(svc *dynamodb.DynamoDB) error {
	_, err := svc.DeleteTable(&dynamodb.DeleteTableInput{
		TableName: aws.String("user"),
	})
	return err
}
