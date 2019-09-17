// +build unit

package mysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	apperror "github.com/jayvib/app/apperr"
	"github.com/jayvib/app/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetByID(t *testing.T) {
	now := time.Now()
	mockArticle := &model.Article{
		ID:        "uniqueid",
		Title:     "Pirate King",
		Content:   "Luffy will be the next Pirate King",
		Author:    &model.Author{ID: "uniqueid"},
		CreatedAt: now,
		UpdatedAt: now,
	}
	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "created_at", "updated_at"}).
		AddRow(mockArticle.ID, mockArticle.Title, mockArticle.Content, mockArticle.Author.ID, mockArticle.CreatedAt, mockArticle.UpdatedAt)

	query := "SELECT id, title, content, author_id, created_at, updated_at FROM article WHERE id = ?"
	t.Run("success", func(t *testing.T) {
		db, smock, err := sqlmock.New()
		require.NoError(t, err)
		smock.
			ExpectQuery(query).
			WithArgs(mockArticle.ID).
			WillReturnRows(rows)

		repo := New(db)
		resArticle, err := repo.GetByID(context.Background(), mockArticle.ID)
		assert.NoError(t, err)
		assert.Equal(t, mockArticle, resArticle)
		smock.ExpectationsWereMet()
	})

	t.Run("failed-not found", func(t *testing.T) {
		db, smock, err := sqlmock.New()
		require.NoError(t, err)
		smock.
			ExpectQuery(query).
			WithArgs("unqieudid").
			WillReturnError(sql.ErrNoRows)

		// I expect to receive an errors.ItemNotFound
		repo := New(db)
		res, err := repo.GetByID(context.Background(), "unqieudid")
		assert.Error(t, err)
		assert.Equal(t, apperror.ItemNotFound, err)
		assert.Nil(t, res)
		smock.ExpectationsWereMet()
	})

	t.Run("unexpected error", func(t *testing.T) {
		db, smock, err := sqlmock.New()
		require.NoError(t, err)
		smock.
			ExpectQuery(query).
			WithArgs(10000).
			WillReturnError(errors.New("unexpected error"))

		u := New(db)
		_, err = u.GetByID(context.Background(), "unqieudi")
		assert.Error(t, err)
		smock.ExpectationsWereMet()
	})
}

func TestGetByTitle(t *testing.T) {
	now := time.Now()
	mockArticle := &model.Article{
		ID:        "uniqueid",
		Title:     "Pirate King",
		Content:   "Luffy will be the next Pirate King",
		Author:    &model.Author{ID: "unqieuid"},
		CreatedAt: now,
		UpdatedAt: now,
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "created_at", "updated_at"}).
		AddRow(mockArticle.ID, mockArticle.Title, mockArticle.Content, mockArticle.Author.ID, mockArticle.CreatedAt, mockArticle.UpdatedAt)
	// Using the table driven test
	query := "SELECT id, title, content, author_id, created_at, updated_at FROM article WHERE title = ?"
	onSuccess := func(t *testing.T) {
		db, smock, err := sqlmock.New()
		require.NoError(t, err)
		smock.
			ExpectQuery(query).
			WithArgs(mockArticle.Title).
			WillReturnRows(rows)

		repo := New(db)
		resArticle, err := repo.GetByTitle(context.Background(), mockArticle.Title)
		assert.NoError(t, err)
		assert.Equal(t, mockArticle, resArticle)
		smock.ExpectationsWereMet()
	}

	notFound := func(t *testing.T) {
		db, smock, err := sqlmock.New()
		require.NoError(t, err)
		smock.
			ExpectQuery(query).
			WithArgs("not exist").
			WillReturnError(sql.ErrNoRows)

		repo := New(db)
		_, err = repo.GetByTitle(context.Background(), "not exist")
		assert.Error(t, err)
		assert.Equal(t, apperror.ItemNotFound, err)
		smock.ExpectationsWereMet()
	}

	unexpectedErr := func(t *testing.T) {
		db, smock, err := sqlmock.New()
		require.NoError(t, err)
		smock.
			ExpectQuery(query).
			WithArgs("not exist").
			WillReturnError(errors.New("unexpected error"))

		repo := New(db)
		_, err = repo.GetByTitle(context.Background(), "not exist")
		assert.Error(t, err)
		smock.ExpectationsWereMet()
	}

	tests := []struct {
		testName string
		testFunc func(*testing.T)
	}{
		{
			"success",
			onSuccess,
		},
		{
			"item not found",
			notFound,
		},
		{
			"unexpected error",
			unexpectedErr,
		},
	}

	for _, d := range tests {
		t.Run(d.testName, d.testFunc)
	}

}

func TestUpdate(t *testing.T) {
	mockArticle := &model.Article{
		ID:      "uniqueid",
		Title:   "Pirate King",
		Content: "Luffy will be the next Pirate King!",
		Author: &model.Author{
			ID:   "unqueid",
			Name: "Luffy Monkey",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	t.Run("success", func(t *testing.T) {
		db, smock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		smock.ExpectBegin()
		smock.ExpectExec("UPDATE article").
			WithArgs(
				mockArticle.Title,
				mockArticle.Content,
				mockArticle.Author.ID,
				mockArticle.UpdatedAt,
				mockArticle.ID,
			).
			WillReturnResult(sqlmock.NewResult(12, 1))
		smock.ExpectCommit()

		u := New(db)
		err = u.Update(context.Background(), mockArticle)
		assert.NoError(t, err)
		assert.NoError(t, smock.ExpectationsWereMet())
	})
}

func TestStore(t *testing.T) {
	mockArticle := &model.Article{
		ID:      "uniqueid",
		Title:   "Pirate King",
		Content: "Luffy will be the next Pirate King!",
		Author: &model.Author{
			ID:   "uniqueid",
			Name: "Luffy Monkey",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		copyArticle := &(*mockArticle)
		copyArticle.ID = "uniqueid"
		db, smock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		smock.ExpectBegin()
		smock.ExpectExec("INSERT INTO article").WithArgs(
			mockArticle.ID,
			mockArticle.Title,
			mockArticle.Content,
			mockArticle.Author.ID,
			mockArticle.CreatedAt,
			mockArticle.UpdatedAt,
		).WillReturnResult(sqlmock.NewResult(12, 1))
		smock.ExpectCommit()
		repo := New(db)
		err = repo.Store(context.Background(), copyArticle)
		require.NoError(t, err)
		assert.Equal(t, mockArticle, copyArticle)
		assert.NoError(t, smock.ExpectationsWereMet())
	})
}

func TestDelete(t *testing.T) {
	mockArticle := &model.Article{
		ID:      "uniqueid",
		Title:   "Pirate King",
		Content: "Luffy will be the next Pirate King!",
		Author: &model.Author{
			ID:   "uniquid",
			Name: "Luffy Monkey",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		db, smock, err := sqlmock.New()
		require.NoError(t, err)

		smock.ExpectBegin()
		smock.ExpectExec("DELETE FROM article").
			WithArgs(mockArticle.ID).
			WillReturnResult(sqlmock.NewResult(12, 1))
		smock.ExpectCommit()
		repo := New(db)
		err = repo.Delete(context.Background(), mockArticle.ID)
		assert.NoError(t, err)
		require.NoError(t, smock.ExpectationsWereMet())
	})

	t.Run("failed-item not found", func(t *testing.T) {
		// TODO: Need to implement
	})
}

func TestFetch(t *testing.T) {
	db, smock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mockArticles := []*model.Article{
		{
			ID:      "uniqiudid",
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
			ID:      "unqieuid",
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

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "updated_at", "created_at"})
	for _, article := range mockArticles {
		rows.AddRow(article.ID, article.Title, article.Content, article.Author.ID, article.UpdatedAt, article.CreatedAt)
	}

	// need to escape the "?" character as per this issue:
	// https://github.com/DATA-DOG/go-sqlmock/issues/70
	query := "SELECT id, title, content, author_id, updated_at, created_at FROM article WHERE created_at > \\? ORDER BY created_at LIMIT \\?"
	smock.ExpectQuery(query).WillReturnRows(rows)
	a := New(db)
	cursor := EncodeCursor(mockArticles[1].CreatedAt)
	num := 2
	list, nextCursor, err := a.Fetch(context.TODO(), cursor, num)
	assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}
