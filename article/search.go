package article

import (
	"context"
	internalsearch "github.com/jayvib/app/internal/app/search"
)

// Use to create a mock object for the SearchEngine
//go:generate mockery --name SearchEngine

// SearchEngine represents a search engine.
type SearchEngine interface {
	// Search takes context and input that will be use for the query.
	// It returns a result of the query and an error if has any.
	Search(ctx context.Context, input internalsearch.Input) (result *internalsearch.Result, err error)
}
