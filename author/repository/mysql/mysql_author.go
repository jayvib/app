package mysql

import (
	"context"
	"database/sql"

	"github.com/jayvib/clean-architecture/apperr"
	"github.com/jayvib/clean-architecture/author"
	"github.com/jayvib/clean-architecture/model"
	"github.com/sirupsen/logrus"
)

func New(db *sql.DB) author.Repository {
	return &mysqlAuthorRepo{
		db: db,
	}
}

type mysqlAuthorRepo struct {
	db *sql.DB
}

func (m *mysqlAuthorRepo) GetByID(ctx context.Context, id string) (*model.Author, error) {
	query := "SELECT id, name, created_at, updated_at FROM author WHERE id = ?"
	return m.getOne(ctx, query, id)
}

func (m *mysqlAuthorRepo) getOne(ctx context.Context, query string, args ...interface{}) (author *model.Author, err error) {
	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer func() {
		err = stmt.Close()
	}()

	row := stmt.QueryRowContext(ctx, args...)
	author = new(model.Author)
	err = row.Scan(
		&author.ID,
		&author.Name,
		&author.CreatedAt,
		&author.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperr.ItemNotFound
		}
		logrus.Error(err)
		return nil, err
	}
	return author, nil
}

func (m *mysqlAuthorRepo) Store(ctx context.Context, u *model.Author) (err error) {
	query := "INSERT INTO author(id, name, created_at, updated_at) VALUES (?, ?, ?, ?)"
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			err = tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(ctx, query, u.ID, u.Name, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
