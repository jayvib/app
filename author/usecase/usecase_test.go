// +build unit

package usecase_test

import (
	"context"
	"github.com/jayvib/clean-architecture/author/mocks"
	"github.com/jayvib/clean-architecture/author/usecase"
	"github.com/jayvib/clean-architecture/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var mockAuthor = &model.Author{
	ID:        "uniqueid",
	Name:      "Luffy Monkey",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

func TestGetByID(t *testing.T) {
	repo := new(mocks.Repository)

	t.Run("success", func(t *testing.T) {
		repo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(mockAuthor, nil).Once()
		uc := usecase.New(repo)
		au, err := uc.GetByID(context.Background(), mockAuthor.ID)
		require.NoError(t, err)
		assert.Equal(t, mockAuthor, au)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		// TODO: Implement me.
		t.SkipNow()
	})
}

func TestStore(t *testing.T) {
	repo := new(mocks.Repository)

	t.Run("success", func(t *testing.T) {
		repo.On("Store", mock.Anything, mock.AnythingOfType("*model.Author")).
			Return(nil).Once()
		uc := usecase.New(repo)
		copyAuthor := &(*mockAuthor)
		copyAuthor.ID = ""
		err := uc.Store(context.Background(), copyAuthor)
		require.NoError(t, err)
		//assert.NotEmpty(t, copyAuthor.ID)
		repo.AssertExpectations(t)
	})
}
