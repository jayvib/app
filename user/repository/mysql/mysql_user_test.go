// +build unit

package mysql

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jayvib/clean-architecture/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *sqlmock.Rows, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	rows := sqlmock.NewRows([]string{
		"id", "firstname", "lastname",
		"email", "username", "password",
		"created_at", "updated_at",
	})
	return db, mock, rows, err

}

func TestMySqlUserRepository(t *testing.T) {
	now := time.Now()
	mockUsers := []*model.User{
		{
			ID:        "uniqueid",
			Firstname: "Luffy",
			Lastname:  "Monkey",
			Email:     "luffy.monkey@test.com",
			Username:  "luffy.monkey",
			Password:  "pirateking",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	t.Run("GetByID", func(t *testing.T) {
		db, mock, rows, err := setup(t)
		require.NoError(t, err)
		defer db.Close()
		repo := New(db)
		for _, u := range mockUsers {
			rows.AddRow(
				u.ID,
				u.Firstname,
				u.Lastname,
				u.Email,
				u.Username,
				u.Password,
				u.CreatedAt,
				u.UpdatedAt,
			)
		}
		// Test the GetById
		query := "SELECT id,firstname,lastname,email,username,password,created_at,updated_at FROM user WHERE id = ?"
		mock.ExpectQuery(query).WillReturnRows(rows)
		u, err := repo.GetByID(context.Background(), "uniqueid")
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, mockUsers[0], u)
	})
	t.Run("GetByUsername", func(t *testing.T) {
		db, mock, rows, err := setup(t)
		require.NoError(t, err)
		defer db.Close()
		repo := New(db)
		for _, u := range mockUsers {
			rows.AddRow(
				u.ID,
				u.Firstname,
				u.Lastname,
				u.Email,
				u.Username,
				u.Password,
				u.CreatedAt,
				u.UpdatedAt,
			)
		}
		// Test the GetById
		query := "SELECT id,firstname,lastname,email,username,password,created_at,updated_at FROM user WHERE username = ?"
		mock.ExpectQuery(query).WillReturnRows(rows)
		username := "luffy.monkey"
		u, err := repo.GetByUsername(context.Background(), username)
		assert.True(t, err == nil, "error found")
		assert.Equal(t, mockUsers[0], u)

	})
	t.Run("GetByEmail", func(t *testing.T) {
		db, mock, rows, err := setup(t)
		require.NoError(t, err)
		defer db.Close()
		repo := New(db)
		for _, u := range mockUsers {
			rows.AddRow(
				u.ID,
				u.Firstname,
				u.Lastname,
				u.Email,
				u.Username,
				u.Password,
				u.CreatedAt,
				u.UpdatedAt,
			)
		}
		query := "SELECT id,firstname,lastname,email,username,password,created_at,updated_at FROM user WHERE email = ?"
		mock.ExpectQuery(query).WillReturnRows(rows)
		email := "luffy.monkey@test.com"
		u, err := repo.GetByEmail(context.Background(), email)
		assert.True(t, err == nil, "error found")
		assert.Equal(t, mockUsers[0], u)
	})
	t.Run("Update", func(t *testing.T) {
		db, mock, _, err := setup(t)
		require.NoError(t, err)
		defer db.Close()
		repo := New(db)
		// Test the GetById
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE user").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = repo.Update(context.Background(), mockUsers[0])
		require.NoError(t, err, "Updating the user received an error")
		err = mock.ExpectationsWereMet()
		require.NoError(t, err, "checking the expectation returns an error")
	})
	t.Run("Store", func(t *testing.T) {
		newUser := &model.User{
			ID:        mockUsers[0].ID,
			Firstname: mockUsers[0].Firstname,
			Lastname:  mockUsers[0].Lastname,
			Email:     mockUsers[0].Email,
			Username:  mockUsers[0].Username,
			Password:  mockUsers[0].Password,
			CreatedAt: mockUsers[0].CreatedAt,
			UpdatedAt: mockUsers[0].UpdatedAt,
		}
		_ = newUser

		db, mock, _, err := setup(t)
		require.NoError(t, err)
		defer db.Close()
		repo := New(db)
		// Test the GetById
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO user").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		err = repo.Store(context.Background(), newUser)
		require.NoError(t, err, "Storing the user received an error")
		//assert.Equal(t, mockUsers[0], newUser, "user not equal")
		assert.NotEmpty(t, newUser.ID)
		err = mock.ExpectationsWereMet()
		require.NoError(t, err, "checking the expectation returns an error")
	})
	t.Run("Delete", func(t *testing.T) {
		//logrus.SetLevel(logrus.DebugLevel)
		db, mock, _, err := setup(t)
		require.NoError(t, err)
		defer db.Close()
		repo := New(db)
		// Test the GetById
		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM user").WithArgs(mockUsers[0].ID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		err = repo.Delete(context.Background(), mockUsers[0].ID)
		require.NoError(t, err, "Deleting the user received an error")
		err = mock.ExpectationsWereMet()
		require.NoError(t, err, "checking the expectation returns an error")

	})
}
