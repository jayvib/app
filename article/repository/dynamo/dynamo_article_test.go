// +build unit

package dynamo

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jayvib/clean-architecture/model"
	"github.com/jayvib/clean-architecture/pkg/mocks"
	apptime "github.com/jayvib/clean-architecture/time"
	"github.com/jayvib/clean-architecture/utils/generateutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

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

func TestFetch(t *testing.T) {
	dynamoMock := new(mocks.DynamoDBAPI)
	output := toScanOuput(t, articles...)
	t.Run("Success", func(t *testing.T) {
		dynamoMock.On("ScanWithContext",
			mock.Anything, mock.AnythingOfType("*dynamodb.ScanInput"), mock.Anything).
			Return(output, nil)

		repo := New(dynamoMock)
		cursor := apptime.EncodeCursor(time.Now().Add(-1 * time.Hour))
		num := 2
		ars, nextCursor, err := repo.Fetch(context.Background(), cursor, num)
		require.NoError(t, err)
		assert.Empty(t, nextCursor)
		assert.Len(t, ars, 2)
	})
}

func TestGetByID(t *testing.T) {
	dynamoMock := new(mocks.DynamoDBAPI)
	av, err := dynamodbattribute.MarshalMap(articles[0])
	require.NoError(t, err)
	output := &dynamodb.GetItemOutput{
		Item: av,
	}
	t.Run("Success", func(t *testing.T) {
		dynamoMock.On("GetItemWithContext",
			mock.Anything, mock.AnythingOfType("*dynamodb.GetItemInput"), mock.Anything).
			Return(output, nil).Once()

		repo := New(dynamoMock)
		artcl, err := repo.GetByID(context.Background(), articles[0].ID)
		require.NoError(t, err)
		assert.Equal(t, articles[0].Title, artcl.Title)
		dynamoMock.AssertExpectations(t)
	})
}

func TestGetByTitle(t *testing.T) {
	dynamoMock := new(mocks.DynamoDBAPI)
	artcle := articles[0]
	itemAv, err := dynamodbattribute.MarshalMap(artcle)
	require.NoError(t, err)
	output := &dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			itemAv,
		},
		Count: aws.Int64(1),
	}

	t.Run("Success", func(t *testing.T) {
		dynamoMock.On("ScanWithContext",
			mock.Anything, mock.AnythingOfType("*dynamodb.ScanInput"), mock.Anything).
			Return(output, nil).Once()

		repo := New(dynamoMock)
		gotArticle, err := repo.GetByTitle(context.Background(), articles[0].Title)
		require.NoError(t, err)
		assert.Equal(t, artcle.ID, gotArticle.ID)
		dynamoMock.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	dbMock := new(mocks.DynamoDBAPI)
	t.Run("Success", func(t *testing.T) {
		article := articles[0]
		dbMock.On("UpdateItemWithContext",
			mock.Anything, mock.AnythingOfType("*dynamodb.UpdateItemInput"), mock.Anything).
			Return(nil, nil).Once()

		repo := New(dbMock)
		err := repo.Update(context.Background(), article)
		require.NoError(t, err)
		dbMock.AssertExpectations(t)
	})
}

func TestStore(t *testing.T) {
	dbMock := new(mocks.DynamoDBAPI)
	t.Run("Success", func(t *testing.T) {
		article := articles[0]
		dbMock.On("PutItemWithContext",
			mock.Anything, mock.AnythingOfType("*dynamodb.PutItemInput"), mock.Anything).
			Return(nil, nil).Once()
		repo := New(dbMock)
		err := repo.Store(context.Background(), article)
		require.NoError(t, err)
		dbMock.AssertExpectations(t)
	})

}

func TestDelete(t *testing.T) {
	dbMock := new(mocks.DynamoDBAPI)
	t.Run("Success", func(t *testing.T) {
		article := articles[0]
		dbMock.On("DeleteItemWithContext",
			mock.Anything, mock.AnythingOfType("*dynamodb.DeleteItemInput"), mock.Anything).
			Return(nil, nil).Once()

		repo := New(dbMock)
		err := repo.Delete(context.Background(), article.ID)
		require.NoError(t, err)
		dbMock.AssertExpectations(t)
	})
}

func toScanOuput(t *testing.T, as ...*model.Article) *dynamodb.ScanOutput {
	output := new(dynamodb.ScanOutput)
	for _, a := range as {
		av, err := dynamodbattribute.MarshalMap(a)
		require.NoError(t, err)
		if output.Items == nil {
			output.Items = []map[string]*dynamodb.AttributeValue{
				av,
			}
		} else {
			output.Items = append(output.Items, av)
		}
	}
	output.Count = aws.Int64(int64(len(as)))
	return output
}
