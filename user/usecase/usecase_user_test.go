// +build unit

package usecase

import (
	"context"
	"github.com/jayvib/clean-architecture/apperr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"

	authormocks "github.com/jayvib/clean-architecture/author/mocks"
	"github.com/jayvib/clean-architecture/model"
	"github.com/jayvib/clean-architecture/user/mocks"
)

var mockUser = &model.User{
	ID:        "uniqueid",
	Firstname: "Luffy",
	Lastname:  "Monkey",
	Username:  "luffy.monkey",
	Email:     "luffy.monkey@gmail.com",
	Password:  "pirateking",
}

func TestGetByID(t *testing.T) {
	mockRepo := new(mocks.Repository)
	authorMockRepo := new(authormocks.Repository)
	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).
			Return(mockUser, nil).Once()
		uc := New(mockRepo, authorMockRepo)
		res, err := uc.GetByID(context.Background(), mockUser.ID)
		require.NoError(t, err)
		assert.Equal(t, mockUser, res)
		mockRepo.AssertExpectations(t)
	})
	t.SkipNow()
	t.Run("not found", func(t *testing.T) {
		// TODO: Implement me!
		t.SkipNow()
	})

	t.Run("unexpected error", func(t *testing.T) {
		// TODO: Implement Me!
		t.SkipNow()
	})
}

func TestGetByEmail(t *testing.T) {
	mockRepo := new(mocks.Repository)
	authorMockRepo := new(authormocks.Repository)
	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).
			Return(mockUser, nil).Once()

		uc := New(mockRepo, authorMockRepo)
		res, err := uc.GetByEmail(context.Background(), mockUser.Email)
		require.NoError(t, err)
		assert.Equal(t, mockUser, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetByUsername(t *testing.T) {
	mockRepo := new(mocks.Repository)
	authorMockRepo := new(authormocks.Repository)
	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetByUsername", mock.Anything, mock.AnythingOfType("string")).
			Return(mockUser, nil).Once()

		uc := New(mockRepo, authorMockRepo)
		res, err := uc.GetByUsername(context.Background(), mockUser.Username)
		require.NoError(t, err)
		assert.Equal(t, mockUser, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestStore(t *testing.T) {
	repo := new(mocks.Repository)
	authorMockRepo := new(authormocks.Repository)
	t.Run("success", func(t *testing.T) {
		repo.On("Store", mock.Anything, mock.AnythingOfType("*model.User")).
			Return(nil).Once()
		repo.On("GetByUsername", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, apperr.ItemNotFound)
		repo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, apperr.ItemNotFound)

		authorMockRepo.On("Store", mock.Anything, mock.AnythingOfType("*model.Author")).
			Return(nil).Once()

		copyUser := *mockUser
		uc := New(repo, authorMockRepo)
		err := uc.Store(context.Background(), &copyUser)
		assert.NotEmpty(t, copyUser.ID)
		assert.NotEqual(t, mockUser.Password, copyUser.Password,
			"current: %s got: %s", mockUser.Password, copyUser.Password)
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	repo := new(mocks.Repository)
	authorMockRepo := new(authormocks.Repository)
	t.Run("success", func(t *testing.T) {
		repo.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).
			Return(nil).Once()
		copyUser := &(*mockUser)
		uc := New(repo, authorMockRepo)
		err := uc.Update(context.Background(), copyUser)
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	repo := new(mocks.Repository)
	authorMockRepo := new(authormocks.Repository)
	t.Run("success", func(t *testing.T) {
		repo.On("Delete", mock.Anything, mock.AnythingOfType("string")).
			Return(nil).Once()
		uc := New(repo, authorMockRepo)
		err := uc.Delete(context.Background(), mockUser.ID)
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})
}
