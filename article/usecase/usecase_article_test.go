// +build unit

package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	myerr "github.com/jayvib/app/apperr"
	"github.com/jayvib/app/article/mocks"
	authorMocks "github.com/jayvib/app/author/mocks"
	"github.com/jayvib/app/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetByID(t *testing.T) {
	mockArticleRepo := new(mocks.Repository)
	mockAuthor := &model.Author{
		ID:   "uniqueid",
		Name: "Luffy Monkey",
	}
	mockArticle := &model.Article{
		Title:   "One Piece",
		Content: "I will be the Pirate King!",
		Author:  mockAuthor,
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(mockArticle, nil).
			Once()

		mockAuthorRepo := new(authorMocks.Repository)
		mockAuthorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(mockAuthor, nil)
		u := New(mockArticleRepo, mockAuthorRepo, time.Second*2)
		a, err := u.GetByID(context.Background(), mockArticle.ID)
		assert.NoError(t, err)
		assert.NotNil(t, a)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})

	t.Run("failed", func(t *testing.T) {
		mockArticleRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, errors.New("Unexpected")).
			Once()

		mockAuthorRepo := new(authorMocks.Repository)
		u := New(mockArticleRepo, mockAuthorRepo, time.Second*2)
		a, err := u.GetByID(context.Background(), mockArticle.ID)
		assert.Error(t, err)
		assert.Nil(t, a)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})
}

func TestStore(t *testing.T) {
	mockArticleRepo := new(mocks.Repository)
	mockArticle := &model.Article{
		Title:   "Pirate King",
		Content: "Luffy will become the new pirate king",
	}

	t.Run("success", func(t *testing.T) {
		tempMockArticle := &(*mockArticle)
		tempMockArticle.ID = ""
		//t.Logf("%#v", tempMockArticle)
		// When the title is not yet exist then
		// Store the article.
		mockArticleRepo.
			On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, myerr.ItemNotFound).
			Once()

		mockArticleRepo.
			On("Store", mock.Anything, mock.AnythingOfType("*model.Article")).
			Return(nil).
			Once()

		mockAuthorRepo := new(authorMocks.Repository)
		u := New(mockArticleRepo, mockAuthorRepo, time.Second*2)
		// Check when the Title and Content has been modified
		err := u.Store(context.Background(), tempMockArticle)
		assert.NoError(t, err)
		assert.Equal(t, mockArticle.Title, tempMockArticle.Title)
		assert.NotEmpty(t, tempMockArticle.ID)
		assert.NotEmpty(t, tempMockArticle.CreatedAt)
		assert.NotEmpty(t, tempMockArticle.UpdatedAt)
		mockArticleRepo.AssertExpectations(t)
	})

	t.Run("existing title", func(t *testing.T) {
		existingArticle := &(*mockArticle)
		//t.Logf("%#v", existingArticle)
		mockArticleRepo.On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).Return(existingArticle, nil).Once()

		mockAuthor := &model.Author{
			ID:   "uniqueid",
			Name: "Luffy Monkey",
		}
		existingArticle.Author = mockAuthor
		mockAuthorRepo := new(authorMocks.Repository)

		u := New(mockArticleRepo, mockAuthorRepo, time.Second*2)
		err := u.Store(context.Background(), existingArticle)
		assert.Error(t, err)
		assert.Equal(t, myerr.ItemExist, err)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockArticleRepo := new(mocks.Repository)
	mockArticle := &model.Article{
		ID:      "uniqueid",
		Title:   "Pirate King",
		Content: "Luffy will become the new pirate king",
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.
			On("Update", mock.Anything, mockArticle).
			Return(nil).
			Once()

		mockAuthorRepo := new(authorMocks.Repository)
		u := New(mockArticleRepo, mockAuthorRepo, time.Second*2)
		err := u.Update(context.Background(), mockArticle)
		assert.NoError(t, err)
		mockArticleRepo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	mockArticleRepo := new(mocks.Repository)
	mockArticle := &model.Article{
		ID:      "uniqueid",
		Title:   "Pirate King",
		Content: "Luff will become the next Pirate King",
	}
	t.Run("success", func(t *testing.T) {
		mockArticleRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(mockArticle, nil).
			Once()

		mockArticleRepo.
			On("Delete", mock.Anything, mock.AnythingOfType("string")).
			Return(nil).
			Once()

		mockAuthorRepo := new(authorMocks.Repository)
		u := New(mockArticleRepo, mockAuthorRepo, time.Second*2)
		err := u.Delete(context.Background(), mockArticle.ID)
		assert.NoError(t, err)
		mockArticleRepo.AssertExpectations(t)
	})

	t.Run("failed not exist", func(t *testing.T) {
		mockArticleRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, myerr.ItemNotFound).
			Once()

		id := "uniqueid"
		mockAuthorRepo := new(authorMocks.Repository)
		u := New(mockArticleRepo, mockAuthorRepo, time.Second*2)
		err := u.Delete(context.Background(), id)
		assert.Error(t, err)
		assert.Equal(t, err, myerr.ItemNotFound)
		mockArticleRepo.AssertExpectations(t)
	})
}

func TestFetch(t *testing.T) {
	mockArticleRepo := new(mocks.Repository)
	mockArticle := &model.Article{
		ID:      "uniqueid",
		Title:   "Pirate King",
		Content: "Luff will become the next Pirate King",
	}

	mockArticles := make([]*model.Article, 0)
	mockArticles = append(mockArticles, mockArticle)

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.
			On("Fetch", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int")).
			Return(mockArticles, "next-cursor", nil).
			Once()

		mockAuthor := &model.Author{
			ID:   "uniqueid",
			Name: "Luffy Monkey",
		}

		mockArticle.Author = mockAuthor

		mockAuthorRepo := new(authorMocks.Repository)
		mockAuthorRepo.
			On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(mockAuthor, nil).
			Once()

		u := New(mockArticleRepo, mockAuthorRepo, time.Second*2)
		num := 1
		cursor := "12"
		expectedCursor := "next-cursor"
		res, nextCursor, err := u.Fetch(context.Background(), cursor, num)
		assert.NoError(t, err)
		assert.Equal(t, expectedCursor, nextCursor)
		assert.Equal(t, mockArticles, res)

		mockArticleRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})

	t.Run("error failed", func(t *testing.T) {
		mockArticleRepo.
			On("Fetch", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int")).
			Return(nil, "", errors.New("unexpected error")).
			Once()

		mockAuthorRepo := new(authorMocks.Repository)
		u := New(mockArticleRepo, mockAuthorRepo, time.Second*2)
		num := 1
		cursor := "12"
		res, nextCursor, err := u.Fetch(context.Background(), cursor, num)
		assert.Error(t, err)
		assert.Empty(t, nextCursor)
		assert.Nil(t, res)
		mockArticleRepo.AssertExpectations(t)
	})

}
