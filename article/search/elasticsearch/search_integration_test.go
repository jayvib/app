// +build integration,elasticsearch

package elasticsearch_test

import (
	"context"
	"encoding/json"
	"github.com/jayvib/app/apperr"
	articlesearches "github.com/jayvib/app/article/search/elasticsearch"
	"github.com/jayvib/app/config"
	internalsearch "github.com/jayvib/app/internal/app/search"
	"github.com/jayvib/app/internal/app/search/elasticsearch/testutil"
	"github.com/jayvib/app/model"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

var client *elastic.Client

const (
	elasticsearchTestURL = "http://localhost:9200"
	index                = "article"
)

// For safety purpose. Do this when doing integration test.
func init() {
	os.Setenv(config.AppEnvironmentKey, config.DevelopmentEnv)
}

func TestSearch(t *testing.T) {
	search := articlesearches.New(client)

	cases := []struct {
		name         string
		ctx          context.Context
		input        internalsearch.Input
		setup        func(t *testing.T)
		assertResult func(t *testing.T, res *internalsearch.Result)
		wantErr      bool
		assertErr    func(t *testing.T, err error)
	}{
		{
			name: "Query in content with total hits is the same with size",
			ctx:  context.Background(),
			input: internalsearch.Input{
				Query: "content=luffy",
				Size:  5,
			},
			setup: func(t *testing.T) {
				testutil.LoadSampleDataFromProvider(t, client, index, provider(t, "search.input"))
			},
			assertResult: func(t *testing.T, res *internalsearch.Result) {
				assert.Len(t, res.Data, 5)
				assert.Equal(t, res.TotalHits, 5)
				assert.Equal(t, 0, res.Next)
			},
			wantErr: false,
		},
		{
			name: "Query with the size is less than the total hits",
			ctx:  context.Background(),
			input: internalsearch.Input{
				Query: "mollit",
				Size:  5,
			},
			setup: func(t *testing.T) {}, // Use the previous data
			assertResult: func(t *testing.T, res *internalsearch.Result) {
				assert.Len(t, res.Data, 5)
				assert.Equal(t, 15, res.TotalHits)
				assert.Equal(t, 5, res.Next)
			},
			wantErr: false,
		},
		{
			name: "Empty query string",
			ctx:  context.Background(),
			input: internalsearch.Input{
				Query: "",
				Size:  5,
			},
			setup:        func(t *testing.T) {},
			assertResult: func(t *testing.T, res *internalsearch.Result) {},
			wantErr:      true,
			assertErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				if aerr, ok := err.(apperr.Error); ok {
					assert.Equal(t, apperr.ValidationErr, aerr.Code())
				} else {
					assert.Fail(t, "expecting apperr.Error type but didn't receive one")
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.setup(t)
			result, err := search.Search(c.ctx, c.input)
			if c.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			c.assertResult(t, result)
		})
	}
}

func TestMain(m *testing.M) {
	var err error
	client, err = testutil.Setup(elasticsearchTestURL, index, "mapping.json")
	if err != nil {
		panic(err)
	}
	exitVal := m.Run()
	os.Exit(exitVal)
}

func provider(t *testing.T, filename string) func() map[string]interface{} {
	file, err := os.Open(filepath.Join("testdata", filename))
	require.NoError(t, err)
	defer func() {
		err = file.Close()
		require.NoError(t, err)
	}()

	var articles []*model.Article
	err = json.NewDecoder(file).Decode(&articles)
	require.NoError(t, err)

	mapArticles := make(map[string]interface{})
	for _, a := range articles {
		mapArticles[a.ID] = a
	}
	return func() map[string]interface{} {
		return mapArticles
	}
}
