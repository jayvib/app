package service

import (
	"context"
	"github.com/jayvib/app/utils/generateutil"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"time"

	apperrors "github.com/jayvib/app/apperr"
	"github.com/jayvib/app/article"
	"github.com/jayvib/app/author"
	"github.com/jayvib/app/model"
)

func New(artr article.Repository, autr author.Repository, duration time.Duration) article.Service {
	return &articleUsecase{
		articleRepo:    artr,
		authorRepo:     autr,
		contextTimeout: duration,
	}
}

type articleUsecase struct {
	articleRepo article.Repository

	// TODO: Instead of incorporating author repository
	// interface use event bus.
	//
	// Read the "Decoupling the components" section.
	// https://herbertograca.com/2017/11/16/explicit-architecture-01-ddd-hexagonal-onion-clean-cqrs-how-i-put-it-all-together/
	authorRepo     author.Repository
	contextTimeout time.Duration
}

func (u *articleUsecase) Fetch(ctx context.Context, cursor string, num int) (ars []*model.Article, nexCursor string, err error) {
	if num == 0 {
		num = 10
	}
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	listArticle, nextCursor, err := u.articleRepo.Fetch(ctx, cursor, num)
	if err != nil {
		logrus.Error(err)
		return nil, "", err
	}

	// Fill the author details
	listArticle, err = u.fillAuthors(ctx, listArticle)
	if err != nil {
		logrus.Error(err)
		return nil, "", err
	}

	return listArticle, nextCursor, nil
}

func (u *articleUsecase) fillAuthors(ctx context.Context, articles []*model.Article) ([]*model.Article, error) {
	g, ctx := errgroup.WithContext(ctx)
	mapAuthors := make(map[string]*model.Author)

	for _, art := range articles {
		if art.Author != nil {
			mapAuthors[art.Author.ID] = new(model.Author)
		}
	}

	chanAuthor := make(chan *model.Author)

	for authorID := range mapAuthors {
		authorID := authorID // variable shadowing
		g.Go(func() error {
			res, err := u.authorRepo.GetByID(ctx, authorID)
			if err != nil {
				return err
			}
			select {
			case <-ctx.Done():
				return nil
			case chanAuthor <- res:
			}
			return nil
		})
	}

	go func() {
		if err := g.Wait(); err != nil {
			logrus.Error(err)
			return
		}
		close(chanAuthor)
	}()

	for auth := range chanAuthor {
		if auth != nil {
			mapAuthors[auth.ID] = auth
		}
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// Fill the author
	for _, item := range articles {
		// check if the author exist
		if item.Author != nil {
			if a, ok := mapAuthors[item.Author.ID]; ok {
				item.Author = a
			}
		}
	}

	return articles, nil
}

func (u *articleUsecase) GetByID(ctx context.Context, id string) (*model.Article, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	art, err := u.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if art.Author != nil {
		auth, err := u.authorRepo.GetByID(ctx, art.Author.ID)
		if err != nil {
			return nil, err
		}
		art.Author = auth
	}

	return art, nil
}

func (u *articleUsecase) GetByTitle(ctx context.Context, title string) (*model.Article, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	res, err := u.articleRepo.GetByTitle(ctx, title)
	if err != nil {
		return nil, err
	}
	if res.Author != nil {
		resAuthor, err := u.authorRepo.GetByID(ctx, res.Author.ID)
		if err != nil {
			return nil, err
		}
		res.Author = resAuthor
	}
	return res, nil
}

func (u *articleUsecase) Update(ctx context.Context, ar *model.Article) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	ar.UpdatedAt = time.Now()
	return u.articleRepo.Update(ctx, ar)
}

func (u *articleUsecase) Store(ctx context.Context, ar *model.Article) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	existedArticle, err := u.articleRepo.GetByTitle(ctx, ar.Title)
	if err != nil && err != apperrors.ItemNotFound {
		return err
	}
	if existedArticle != nil {
		return apperrors.ItemExist
	}

	ar.ID = generateutil.GenerateID(ar.TableName())
	ar.CreatedAt = time.Now()
	ar.UpdatedAt = time.Now()
	err = u.articleRepo.Store(ctx, ar)
	if err != nil {
		return err
	}

	return nil
}

func (u *articleUsecase) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	res, err := u.articleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if res == nil {
		return apperrors.ItemNotFound
	}

	return u.articleRepo.Delete(ctx, id)
}
