package user

import (
	"context"

	"github.com/jayvib/clean-architecture/model"
)

//go:generate mockery -name=Repository

type Repository interface {
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Store(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
}
