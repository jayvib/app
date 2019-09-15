package mysql

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	apperror "github.com/jayvib/clean-architecture/apperr"
	"github.com/jayvib/clean-architecture/article"
	"github.com/jayvib/clean-architecture/model"
	"github.com/sirupsen/logrus"
	"time"
)

var timeFormat = "2006-01-02T15:04:05.999Z07:00"

func New(db *sql.DB) article.Repository {
	return &mysqlArticleRepository{db: db}
}

type mysqlArticleRepository struct {
	db *sql.DB
}

func (a *mysqlArticleRepository) Fetch(ctx context.Context, cursor string, num int) ([]*model.Article, string, error) {
	query := `SELECT 
				id, title, content, author_id, updated_at, created_at
			  FROM 
				article 
		      WHERE 
				created_at > ? 
			  ORDER BY 
				created_at 
			  LIMIT ?`

	decodedCursor, err := DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", apperror.BadParamInput
	}

	res, err := a.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	nextCursor := ""
	if len(res) == int(num) {
		nextCursor = EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return res, nextCursor, err
}
func (a *mysqlArticleRepository) GetByID(ctx context.Context, id string) (*model.Article, error) {
	query := "SELECT id, title, content, author_id, created_at, updated_at FROM article WHERE id = ?"
	res, err := a.fetch(ctx, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.ItemNotFound
		}
		return nil, err
	}
	return res[0], nil
}
func (a *mysqlArticleRepository) fetch(ctx context.Context, query string, args ...interface{}) (ars []*model.Article, err error) {
	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		rows.Close()
	}()

	articles := make([]*model.Article, 0)

	for rows.Next() {
		var art model.Article
		var authorId string
		err = rows.Scan(
			&art.ID,
			&art.Title,
			&art.Content,
			&authorId,
			&art.CreatedAt,
			&art.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		art.Author = &model.Author{ID: authorId}
		articles = append(articles, &art)
	}

	if err := rows.Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Debug()
		return nil, err
	}

	if len(articles) == 0 {
		return nil, apperror.ItemNotFound
	}

	return articles, nil
}
func (a *mysqlArticleRepository) GetByTitle(ctx context.Context, title string) (*model.Article, error) {
	query := "SELECT id, title, content, author_id, created_at, updated_at FROM article WHERE title = ?"
	res, err := a.fetch(ctx, query, title)
	logrus.WithFields(logrus.Fields{
		"error": err,
		"title": title,
	}).Debug("Getting by title")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.ItemNotFound
		}
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Debug()
		return nil, err
	}

	if len(res) == 0 {
		return nil, apperror.ItemNotFound
	}

	return res[0], nil
}
func (a *mysqlArticleRepository) Update(ctx context.Context, article *model.Article) (err error) {
	query := "UPDATE article SET title=?, content=?, author_id=?, updated_at=? WHERE id=?"
	tx, err := a.db.Begin()
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

	result, err := tx.ExecContext(ctx, query,
		article.Title, article.Content, article.Author.ID, article.UpdatedAt,
		article.ID)
	if err != nil {
		return err
	}

	// check the affected row
	rowAffectedCount, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffectedCount != 1 {
		return errors.New("rows affected is more than 1")
	}

	return nil
}
func (a *mysqlArticleRepository) Store(ctx context.Context, article *model.Article) (err error) {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		switch err {
		case nil:
			tx.Commit()
		default:
			tx.Rollback()
		}
	}()
	query := "INSERT INTO article(id, title, content, author_id, created_at, updated_at) VALUES(?,?,?,?,?,?)"
	_, err = tx.ExecContext(ctx, query,
		article.ID,
		article.Title,
		article.Content,
		article.Author.ID,
		article.CreatedAt,
		article.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return
}
func (a *mysqlArticleRepository) Delete(ctx context.Context, id string) (err error) {
	tx, err := a.db.Begin()
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

	query := "DELETE FROM article WHERE id=?"
	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return errors.New("rows affected is more than 1")
	}

	return nil
}

func EncodeCursor(t time.Time) string {
	timeString := t.Format(timeFormat)
	return base64.StdEncoding.EncodeToString([]byte(timeString))
}

func DecodeCursor(encodedTime string) (time.Time, error) {
	b, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return time.Time{}, err
	}
	timeString := string(b)
	return time.Parse(timeFormat, timeString)
}
