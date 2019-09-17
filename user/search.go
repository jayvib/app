package user

import (
	"context"
	"github.com/jayvib/app/model"
)

//go:generate mockery -name=SearchEngine

type SearchEngine interface {
	Update(ctx context.Context, user *model.User) error
	Store(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
	SearchByName(ctx context.Context, name string, num int, size int) (users []*model.User, nextPage int, err error)

	// TODO: Use the search input and result in package internal/search.
	Search(ctx context.Context, input SearchInput) (result *SearchResult, err error)
}

type SearchInput struct {
	ScrollId string `json:"scrollId,omitempty"`
	Query    string `json:"query,omitempty" validate:"required"`
	Size     int    `json:"size,omitempty"`
	From     int    `json:"from,omitempty"`
}

type SearchResult struct {
	Users        []*model.User `json:"users"`
	TotalHits    int           `json:"totalHits,omitempty"`
	TookInMillis int           `json:"took,omitempty"`
	ScrollId     string        `json:"scrollId,omitempty"`
	TimedOut     bool          `json:"TimedOut,omitempty"`
	Error        *ErrorDetails `json:"error,omitempty"`
	Next         int           `json:"next,omitempty"`
}

type ErrorDetails struct {
	Type         string                   `json:"type"`
	Reason       string                   `json:"reason"`
	ResourceType string                   `json:"resource.type,omitempty"`
	ResourceId   string                   `json:"resource.id,omitempty"`
	Index        string                   `json:"index,omitempty"`
	Phase        string                   `json:"phase,omitempty"`
	Grouped      bool                     `json:"grouped,omitempty"`
	CausedBy     map[string]interface{}   `json:"caused_by,omitempty"`
	RootCause    []*ErrorDetails          `json:"root_cause,omitempty"`
	FailedShards []map[string]interface{} `json:"failed_shards,omitempty"`
}

type ShardsInfo struct {
	Total      int             `json:"total"`
	Successful int             `json:"successful"`
	Failed     int             `json:"failed"`
	Failures   []*ShardFailure `json:"failures,omitempty"`
}

type ShardFailure struct {
	Index   string                 `json:"_index,omitempty"`
	Shard   int                    `json:"_shard,omitempty"`
	Node    string                 `json:"_node,omitempty"`
	Reason  map[string]interface{} `json:"reason,omitempty"`
	Status  string                 `json:"status,omitempty"`
	Primary bool                   `json:"primary,omitempty"`
}
