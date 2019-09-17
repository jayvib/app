package usecase

import (
	"context"
	"github.com/jayvib/app/apperr"
	"github.com/jayvib/app/author"
	"github.com/jayvib/app/model"
	"github.com/jayvib/app/user"
	"github.com/jayvib/app/utils/crypto"
	"github.com/jayvib/app/utils/generateutil"
	"github.com/sirupsen/logrus"
	"time"
)

func New(repo user.Repository, authrepo author.Repository) user.Usecase {
	return &User{
		repo:       repo,
		authorRepo: authrepo,
	}
}

type User struct {
	repo       user.Repository
	authorRepo author.Repository
}

func (u *User) GetByID(ctx context.Context, id string) (*model.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return u.repo.GetByID(ctx, id)
}

func (u *User) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return u.repo.GetByEmail(ctx, email)
}

func (u *User) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return u.repo.GetByUsername(ctx, username)
}

func (u *User) Update(ctx context.Context, user *model.User) error {
	return u.repo.Update(ctx, user)
}

func (u *User) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func (u *User) Store(ctx context.Context, usr *model.User) error {
	if ctx == nil {
		ctx = context.Background()
	}
	_, err := u.GetByUsername(ctx, usr.Username)
	if err != nil && err != apperr.ItemNotFound {
		return err
	}
	if err == nil {
		return apperr.UsernameAlreadyExist
	}

	res, err := u.GetByEmail(ctx, usr.Email)
	if err != nil && err != apperr.ItemNotFound {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"user":  res,
		"error": err,
	}).Debug()

	if err == nil && res == nil {
		return apperr.EmailAlreadyExist
	}

	usr.CreatedAt = time.Now()
	usr.UpdatedAt = time.Now()
	usr.ID = generateutil.GenerateID(usr.TableName())
	encryptedPass, err := crypto.EncryptPassword(usr.Password)
	if err != nil {
		return err
	}

	usr.Password = encryptedPass

	// Store also to the author repository
	au := &model.Author{
		ID:        usr.ID,
		Name:      usr.FullName(),
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
	}

	logrus.WithFields(logrus.Fields{
		"author_id": au.ID,
	}).Debug()
	err = u.authorRepo.Store(ctx, au)
	if err != nil {
		return err
	}

	return u.repo.Store(ctx, usr)
}
