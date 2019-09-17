package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jayvib/app/apperr"
	"github.com/jayvib/app/model"
	"github.com/jayvib/app/user"
	"github.com/sirupsen/logrus"
)

func New(conn *sql.DB) user.Repository {
	return &mysqlUserRepository{conn: conn}
}

type mysqlUserRepository struct {
	conn *sql.DB
}

func (m *mysqlUserRepository) GetByID(ctx context.Context, id string) (u *model.User, err error) {
	query := "SELECT id,firstname,lastname,email,username,password,created_at,updated_at FROM user WHERE id = ?"
	row, err := m.conn.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = row.Close()
	}()

	users := make([]*model.User, 0)
	for row.Next() {
		var user model.User
		err = row.Scan(
			&user.ID,
			&user.Firstname,
			&user.Lastname,
			&user.Email,
			&user.Username,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err = row.Err(); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if len(users) == 0 {
		return nil, apperr.ItemNotFound
	}
	return users[0], nil
}
func (m *mysqlUserRepository) GetByEmail(ctx context.Context, email string) (u *model.User, err error) {
	query := "SELECT id,firstname,lastname,email,username,password,created_at,updated_at FROM user WHERE email = ?"
	row, err := m.conn.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}
	defer func() {
		row.Close()
	}()
	users := make([]*model.User, 0)
	for row.Next() {
		var user model.User
		row.Scan(
			&user.ID,
			&user.Firstname,
			&user.Lastname,
			&user.Email,
			&user.Username,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err = row.Err(); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if len(users) == 0 {
		return nil, apperr.ItemNotFound
	}
	logrus.WithFields(logrus.Fields{
		"user": users,
	}).Debug()
	return users[0], nil
}
func (m *mysqlUserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query := "SELECT id,firstname,lastname,email,username,password,created_at,updated_at FROM user WHERE username = ?"
	row, err := m.conn.QueryContext(ctx, query, username)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	users := make([]*model.User, 0)
	for row.Next() {
		var user model.User
		row.Scan(
			&user.ID,
			&user.Firstname,
			&user.Lastname,
			&user.Email,
			&user.Username,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err = row.Err(); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if len(users) == 0 {
		return nil, apperr.ItemNotFound
	}
	return users[0], nil
}

func (m *mysqlUserRepository) Update(ctx context.Context, user *model.User) (err error) {

	// use transaction
	tx, err := m.conn.Begin()
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}

	}()

	query := "UPDATE user SET firstname=?, lastname=?, email=?, username=?, password=?, updated_at=? WHERE id = ?"
	_, err = tx.ExecContext(
		ctx, query,
		user.Firstname, user.Lastname, user.Email, user.Username, user.Password, user.UpdatedAt,
		user.ID)
	if err != nil {
		return
	}
	return
}

func (m *mysqlUserRepository) Store(ctx context.Context, user *model.User) (err error) {
	// TODO: This need to check if the email and username is already exist.

	tx, err := m.conn.Begin()
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	query := "INSERT INTO user(id, firstname, lastname, email, username, password, created_at, updated_at) VALUES (?,?,?,?,?,?,?,?)"
	_, err = tx.ExecContext(
		ctx, query, user.ID, user.Firstname,
		user.Lastname, user.Email, user.Username,
		user.Password, user.CreatedAt, user.UpdatedAt)

	return err
}

func (m *mysqlUserRepository) Delete(ctx context.Context, id string) error {
	tx, err := m.conn.Begin()
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	query := "DELETE FROM user WHERE id = ?"
	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	logrus.Debugf("Affected row: %d\n", affected)
	if affected != 1 {
		return fmt.Errorf("Weird! Total Affected Rows: %d\n", affected)
	}

	return err
}
