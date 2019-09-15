package author

import (
	"context"
	"github.com/jayvib/clean-architecture/model"
)

//go:generate mockery -name=Usecase

type Usecase interface {
	GetByID(ctx context.Context, id string) (*model.Author, error)
	Store(ctx context.Context, u *model.Author) error
}
