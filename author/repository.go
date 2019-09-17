package author

import (
	"context"

	"github.com/jayvib/app/model"
)

//go:generate mockery -name=Repository

type Repository interface {
	GetByID(ctx context.Context, id string) (*model.Author, error)
	Store(ctx context.Context, u *model.Author) error
}
