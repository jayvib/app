package service

import (
	"context"
	"github.com/jayvib/app/author"
	"github.com/jayvib/app/model"
)

func New(repo author.Repository) author.Service {
	return &usecase{
		repo: repo,
	}
}

type usecase struct {
	repo author.Repository
}

func (u *usecase) GetByID(ctx context.Context, id string) (*model.Author, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return u.repo.GetByID(ctx, id)
}

func (u *usecase) Store(ctx context.Context, a *model.Author) error {
	// Currently reuse the user id and use it for author id
	// Check first if the usernam
	if ctx == nil {
		ctx = context.Background()
	}
	//a.ID = generateutil.GenerateID(a.TableName())
	return u.repo.Store(ctx, a)
}
