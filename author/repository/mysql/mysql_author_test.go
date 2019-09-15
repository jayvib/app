// +build unit

package mysql

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jayvib/clean-architecture/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var expectedAuthor = &model.Author{
	ID:        "uniqueid",
	Name:      "Luffy Monkey",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

func setup() (*sql.DB, sqlmock.Sqlmock, *sqlmock.Rows, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, nil, err
	}
	rows := sqlmock.NewRows([]string{
		"id", "name", "created_at", "updated_at",
	})
	return db, mock, rows, nil
}

func TestMysqlGetByID(t *testing.T) {
	db, mock, rows, err := setup()
	require.NoError(t, err, "error while setup")
	repo := New(db)

	t.Run("Prepare Statement Error", func(t *testing.T) {
		cdb, _, _ := sqlmock.New()
		cdb.Close()
		crepo := New(cdb)
		_, err := crepo.GetByID(context.Background(), "uniqueid")
		assert.Error(t, err, "expecting an error")
	})

	rows.AddRow(
		expectedAuthor.ID,
		expectedAuthor.Name,
		expectedAuthor.CreatedAt,
		expectedAuthor.UpdatedAt,
	)
	query := "SELECT id, name, created_at, updated_at FROM author WHERE id = ?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(expectedAuthor.ID).WillReturnRows(rows)
	a, err := repo.GetByID(context.Background(), "uniqueid")
	require.NoError(t, err)
	require.NotNil(t, a)
	assert.Equal(t, expectedAuthor, a, "author got not matched with the expected one")
	require.NoError(t, mock.ExpectationsWereMet(), "expectations weren't met")
}

func TestStore(t *testing.T) {
	db, smock, _, err := setup()
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		smock.ExpectBegin()
		smock.ExpectExec("INSERT INTO author").WillReturnResult(sqlmock.NewResult(1, 1))
		smock.ExpectCommit()
		repo := New(db)
		err := repo.Store(context.Background(), expectedAuthor)
		require.NoError(t, err)
		assert.NoError(t, smock.ExpectationsWereMet())
	})
}
