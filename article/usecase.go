package article

import (
	"context"

	"github.com/jayvib/app/model"
)

//go:generate mockery -name=Usecase

type Usecase interface {
	Fetch(ctx context.Context, cursor string, num int) (ars []*model.Article, nexCursor string, err error)
	GetByID(ctx context.Context, id string) (ar *model.Article, err error)
	GetByTitle(ctx context.Context, title string) (ar *model.Article, err error)
	Update(ctx context.Context, ar *model.Article) error
	Store(ctx context.Context, ar *model.Article) error
	Delete(ctx context.Context, id string) error
}
