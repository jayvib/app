// +build unit

package dynamo

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jayvib/clean-architecture/model"
	"github.com/jayvib/clean-architecture/user/mocks"
	"github.com/jayvib/clean-architecture/utils/generateutil"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

var mockUser = &model.User{
	ID:        generateutil.GenerateID("user"),
	Firstname: "Luffy",
	Lastname:  "Monkey",
	Username:  "luffy.monkey",
	Email:     "luffy.monkey@onepice.com",
	Password:  "pirateking",
}

func TestGetByID(t *testing.T) {
	dbmock := new(mocks.DynamoDBAPI)
	itemav, err := MarshalUserToAttributeValue(mockUser)
	//t.Logf("%#v", itemav)
	require.NoError(t, err)
	expected := &dynamodb.GetItemOutput{
		Item: itemav,
	}
	t.Run("Success", func(t *testing.T) {
		dbmock.On("GetItemWithContext", mock.Anything, mock.AnythingOfType("*dynamodb.GetItemInput"), mock.Anything).
			Return(expected, nil).Once()
		repo := New(dbmock)
		u, err := repo.GetByID(context.Background(), mockUser.ID)
		require.NoError(t, err)
		assert.Equal(t, mockUser, u)
		dbmock.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		t.SkipNow()
	})
}

func TestGetByEmail(t *testing.T) {
	dbmock := new(mocks.DynamoDBAPI)
	itemav, err := MarshalUserToAttributeValue(mockUser)
	require.NoError(t, err)
	expected := &dynamodb.ScanOutput{
		Count: aws.Int64(1),
		Items: []map[string]*dynamodb.AttributeValue{
			itemav,
		},
	}
	t.Run("Success", func(t *testing.T) {
		dbmock.On("ScanWithContext",
			mock.Anything, mock.AnythingOfType("*dynamodb.ScanInput"), mock.Anything).
			Return(expected, nil)

		repo := New(dbmock)
		u, err := repo.GetByEmail(context.Background(), mockUser.Email)
		require.NoError(t, err)
		assert.Equal(t, mockUser, u)
		dbmock.AssertExpectations(t)
	})
}

func TestGetByUsername(t *testing.T) {
	dbmock := new(mocks.DynamoDBAPI)
	itemav, err := MarshalUserToAttributeValue(mockUser)
	require.NoError(t, err)
	expected := &dynamodb.ScanOutput{
		Count: aws.Int64(1),
		Items: []map[string]*dynamodb.AttributeValue{
			itemav,
		},
	}
	t.Run("Success", func(t *testing.T) {
		dbmock.On("ScanWithContext",
			mock.Anything, mock.AnythingOfType("*dynamodb.ScanInput"), mock.Anything).
			Return(expected, nil)

		repo := New(dbmock)
		u, err := repo.GetByUsername(context.Background(), mockUser.Username)
		require.NoError(t, err)
		assert.Equal(t, mockUser, u)
		dbmock.AssertExpectations(t)
	})
}

func TestStore(t *testing.T) {
	dbmock := new(mocks.DynamoDBAPI)
	copyUser := mockUser
	t.Run("success", func(t *testing.T) {
		dbmock.On("PutItemWithContext",
			mock.Anything, mock.AnythingOfType("*dynamodb.PutItemInput"), mock.Anything).
			Return(nil, nil)
		repo := New(dbmock)
		err := repo.Store(context.Background(), copyUser)
		require.NoError(t, err)
		dbmock.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	dbmock := new(mocks.DynamoDBAPI)
	t.Run("success", func(t *testing.T) {
		dbmock.On("UpdateItemWithContext",
			mock.Anything, mock.AnythingOfType("*dynamodb.UpdateItemInput"), mock.Anything).
			Return(nil, nil)

		repo := New(dbmock)
		err := repo.Update(context.Background(), mockUser)
		require.NoError(t, err)
		dbmock.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	dbmock := new(mocks.DynamoDBAPI)
	t.Run("success", func(t *testing.T) {
		dbmock.On("DeleteItemWithContext",
			mock.Anything, mock.AnythingOfType("*dynamodb.DeleteItemInput"), mock.Anything).
			Return(nil, nil)

		repo := New(dbmock)
		err := repo.Delete(context.Background(), "uniqueid")
		require.NoError(t, err)
		dbmock.AssertExpectations(t)
	})
}
