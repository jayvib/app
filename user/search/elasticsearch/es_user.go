package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jayvib/app/apperr"
	"github.com/jayvib/app/log"
	"github.com/jayvib/app/model"
	"github.com/jayvib/app/pkg/validator"
	"github.com/jayvib/app/user"
	"github.com/olivere/elastic/v7"
	"reflect"
	"strings"
)

const defaultResultSize = 10

var defaultValidator = validator.New()

// user.SearchEngine implementation reference
var _ user.SearchEngine = (*Search)(nil)

// Opts functions that are available to use.
//
// Please be informed that some nil value params
// should be replaced by the actual object
// implementation. I just use nil value in order to
//  make the compiler happy. 〜(￣▽￣〜)
var _ Opt = SetValidatorOpt(nil)
var _ Opt = SetLoggerOpt(nil)

type Opt func(search *Search)

func New(client *elastic.Client, opts ...Opt) *Search {
	s := &Search{
		client:    client,
		validator: defaultValidator,
	}

	for _, opt := range opts {
		opt(s)
	}
	return s
}

type Search struct {
	client    *elastic.Client // hard to mock, just do integration test directly. (ㄒoㄒ)
	validator validator.Validator
	logger    log.Logger
}

// GetByID takes context in order to allow the caller to stop or cancel the work for this method.
// It also takes an id which should not be and empty string unless it will return an
// apperr.BadParameter error code.
func (s *Search) GetByID(ctx context.Context, id string) (user *model.User, err error) {
	if id == "" {
		return nil, apperr.New(apperr.BadParameter, "id is an empty string", nil)
	}

	res, err := s.client.Get().Index(user.TableName()).Id(id).Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			return nil, apperr.New(apperr.NoItemFound,
				fmt.Sprintf("document with id %s not found", id), err)
		}
	}
	user = new(model.User)
	err = json.Unmarshal(res.Source, user)
	if err != nil {
		return nil, apperr.New(apperr.InternalError, err.Error(), err)
	}

	return user, nil
}

func (s *Search) Update(ctx context.Context, user *model.User) error {
	if err := s.validator.Struct(user); err != nil {
		if validator.IsValidationErr(err) {
			return apperr.New(apperr.ValidationErr, "error while validating user", err)
		}
		return apperr.New(apperr.InternalError, err.Error(), err)
	}
	_, err := s.client.Update().Index(user.TableName()).Id(user.ID).Doc(user).Do(context.Background())
	if err != nil {
		return apperr.New(apperr.InternalError, err.Error(), err)
	}
	return nil
}

func (s *Search) Store(ctx context.Context, user *model.User) error {
	_, err := s.client.Index().Index(user.TableName()).Id(user.ID).BodyJson(user).Do(ctx)
	if err != nil {
		return apperr.New(apperr.InternalError, err.Error(), err)
	}
	return nil
}

func (s *Search) Delete(ctx context.Context, id string) error {
	_, err := s.client.Delete().Index(model.GetUserTableName()).Id(id).Do(ctx)
	if err != nil {
		return apperr.New(apperr.InternalError, err.Error(), err)
	}
	return nil
}

func (s *Search) SearchByName(ctx context.Context, name string, from int, size int) (users []*model.User, next int, err error) {
	if size == 0 {
		size = defaultResultSize
	}
	name = strings.ToLower(name)
	query := elastic.NewTermQuery("firstname", name)
	result, err := s.client.Search().
		Index(model.GetUserTableName()).
		Query(query).
		From(from).Size(size).
		Do(ctx)
	if err != nil {
		return nil, 0, apperr.New(apperr.InternalError, err.Error(), err)
	}
	users = make([]*model.User, 0)
	var userType model.User
	for _, item := range result.Each(reflect.TypeOf(userType)) {
		if u, ok := item.(model.User); ok {
			users = append(users, &u)
		}
	}
	next = from + len(users)
	if next >= int(result.TotalHits()) {
		next = 0
	}

	return users, next, nil
}

// Search is a general query that use the uri parameter to do
// elasticsearch query.
func (s *Search) Search(ctx context.Context, input user.SearchInput) (result *user.SearchResult, err error) {
	// validate the input
	if err := s.validator.Struct(input); err != nil {
		if validator.IsValidationErr(err) {
			aerr := apperr.New(apperr.ValidationErr, "search input validation error", nil)
			return nil, aerr
		}
		return nil, err
	}

	if input.Size == 0 {
		input.Size = defaultResultSize
	}

	log.Infof("query %s size %d from %d", input.Query, input.Size, input.From)

	query := elastic.NewMultiMatchQuery(input.Query)
	res, err := s.client.Search().
		Index(model.GetUserTableName()).
		Query(query).
		From(input.From).Size(input.Size).
		Do(ctx)

	if err != nil {
		if elastic.IsNotFound(err) {
			return nil, apperr.New(apperr.NoItemFound, "No item matched with the query", err)
		}
		return nil, apperr.New(apperr.InternalError, err.Error(), nil)
	}

	users := make([]*model.User, 0)

	for _, item := range res.Each(reflect.TypeOf(model.User{})) {
		if u, ok := item.(model.User); ok {
			users = append(users, &u)
		}
	}

	next := input.From + len(users)
	if next >= int(res.TotalHits()) {
		next = 0
	}
	log.Infof("%#v", users)
	result = &user.SearchResult{
		Users:        users,
		TotalHits:    int(res.TotalHits()),
		TookInMillis: int(res.TookInMillis),
		TimedOut:     res.TimedOut,
		ScrollId:     res.ScrollId,
		Next:         next,
	}

	return result, nil
}

func (s *Search) setValidator(v validator.Validator) {
	s.validator = v
}

func (s *Search) setLogger(l log.Logger) {
	s.logger = l
}

func toUserErrorDetails(err *elastic.ErrorDetails) *user.ErrorDetails {
	return &user.ErrorDetails{
		Type:         err.Type,
		Reason:       err.Reason,
		ResourceType: err.ResourceType,
		ResourceId:   err.ResourceId,
		Index:        err.Index,
		Phase:        err.Phase,
		Grouped:      err.Grouped,
		CausedBy:     err.CausedBy,
		FailedShards: err.FailedShards,
	}
}
