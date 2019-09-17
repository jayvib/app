// +build unit

package dynamo

import (
	"context"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jayvib/app/model"
	"github.com/jayvib/app/pkg/mocks"
	"github.com/jayvib/app/utils/generateutil"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRepository_GetByID(t *testing.T) {
	authorMock := getAuthorMock()
	av, err := dynamodbattribute.MarshalMap(authorMock)
	require.NoError(t, err)
	getItemOutput := &dynamodb.GetItemOutput{
		Item: av,
	}

	authorMock.ID = generateutil.GenerateID(authorMock.TableName())
	svcMock := new(mocks.DynamoDBAPI)
	t.Run("Success", func(t *testing.T) {
		svcMock.On("GetItemWithContext",
			mock.Anything, mock.AnythingOfType("*dynamodb.GetItemInput"), mock.Anything).
			Return(getItemOutput, nil).Once()

		repo := New(svcMock)
		res, err := repo.GetByID(context.Background(), authorMock.ID)
		require.NoError(t, err)
		assert.Equal(t, authorMock.Name, res.Name)
		svcMock.AssertExpectations(t)
	})
}

func TestRepository_Store(t *testing.T) {
	authorMock := getAuthorMock()
	svc := new(mocks.DynamoDBAPI)
	t.Run("Success", func(t *testing.T) {
		svc.On("PutItemWithContext",
			mock.Anything, mock.AnythingOfType("*dynamodb.PutItemInput"), mock.Anything).
			Return(nil, nil)
		repo := New(svc)
		err := repo.Store(context.Background(), authorMock)
		require.NoError(t, err)
		svc.AssertExpectations(t)
	})
}

func getAuthorMock() *model.Author {
	authorMock := &model.Author{
		Name:      "Luffy Monkey",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	authorMock.ID = generateutil.GenerateID(authorMock.TableName())
	return authorMock
}
