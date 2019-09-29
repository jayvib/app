package elasticsearch

import (
	"context"
	"github.com/jayvib/app/apperr"
	internalsearch "github.com/jayvib/app/internal/app/search"
	"github.com/jayvib/app/log"
	"github.com/jayvib/app/model"
	"github.com/jayvib/app/pkg/validator"
	"github.com/olivere/elastic/v7"
	"reflect"
)

const (
	defaultResultSize = 10
)

var (
	defaultValidator = validator.New()
	defaultLogger    = log.NewStandardOutLogger()
)

type Opt func(search *Search)

var _ Opt = SetValidatorOpt(nil)
var _ Opt = SetLoggerOpt(nil)

func New(client *elastic.Client, opts ...Opt) *Search {
	s := &Search{
		validator: defaultValidator,
		logger:    defaultLogger,
		client:    client,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type Search struct {
	validator validator.Validator
	logger    log.Logger
	client    *elastic.Client
}

func (s *Search) Store(ctx context.Context, a *model.Article) error {
	_, err := s.client.Index().Index(model.GetArticleTableName()).Id(a.ID).BodyJson(a).Do(ctx)
	if err != nil {
		return apperr.New(apperr.InternalError, "error while indexing", err)
	}
	return nil
}

func (s *Search) Search(ctx context.Context, input internalsearch.Input) (result *internalsearch.Result, err error) {
	if err := s.validator.Struct(input); err != nil {
		if validator.IsValidationErr(err) {
			aerr := apperr.New(apperr.ValidationErr, "search input validation error", err)
			return nil, aerr
		}
		return nil, apperr.New(apperr.InternalError, err.Error(), err)
	}

	if input.Size == 0 {
		input.Size = defaultResultSize
	}

	query := elastic.NewMultiMatchQuery(input.Query)
	res, err := s.client.Search().
		Index(model.GetArticleTableName()).
		Query(query).
		From(input.From).Size(input.Size).
		Do(ctx)
	if err != nil {
		return nil, apperr.New(apperr.InternalError, err.Error(), err)
	}

	articles := make([]interface{}, 0)
	for _, item := range res.Each(reflect.TypeOf(model.Article{})) {
		if a, ok := item.(model.Article); ok {
			articles = append(articles, &a)
		}
	}

	next := input.From + len(articles)
	if next >= int(res.TotalHits()) {
		next = 0
	}
	result = &internalsearch.Result{
		Data:         articles,
		TotalHits:    int(res.TotalHits()),
		TookInMillis: int(res.TookInMillis),
		TimedOut:     res.TimedOut,
		ScrollId:     res.ScrollId,
		Next:         next,
	}
	return result, nil
}

func (s *Search) SetValidator(v validator.Validator) {
	s.validator = v
}

func (s *Search) SetLogger(l log.Logger) {
	s.logger = l
}

// Options
func SetValidatorOpt(v validator.Validator) Opt {
	return func(search *Search) {
		search.SetValidator(v)
	}
}

func SetLoggerOpt(l log.Logger) Opt {
	return func(search *Search) {
		search.SetLogger(l)
	}
}
