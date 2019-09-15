// +build integration,elasticsearch

package elasticsearch_test

import (
	"context"
	"github.com/jayvib/clean-architecture/pkg/elasticsearch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateIndex(t *testing.T) {
	indexName := "user"
	tdown := setup(t, indexName)
	defer tdown(t)
	t.Run("mapping on file", func(t *testing.T) {
		mappingFileName := "user.json"
		file, err := os.Open(filepath.Join("testdata", mappingFileName))
		require.NoError(t, err)
		result, err := elasticsearch.CreateIndex(indexName, file)
		require.NoError(t, err)
		assert.True(t, result.Acknowledged)
		assert.True(t, result.ShardsAcknowledged)
		assert.Equal(t, indexName, result.Index)
	})

	t.Run("not a json format", func(t *testing.T) {
		strReader := strings.NewReader("thisisnot a json format")
		result, err := elasticsearch.CreateIndex("any", strReader)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("not acknowledged", func(t *testing.T) {
		t.SkipNow()
	})

	t.Run("empty index name", func(t *testing.T) {
		_, err := elasticsearch.CreateIndex("", strings.NewReader("luffy the pirate king"))
		assert.Error(t, err)
		assert.Equal(t, elasticsearch.EmptyIndexNameErr, err)
	})

	t.Run("nil Reader", func(t *testing.T) {
		_, err := elasticsearch.CreateIndex("OncePiece", nil)
		assert.Error(t, err)
		assert.Equal(t, elasticsearch.NilReaderErr, err)
	})
}

func setup(t *testing.T, indexName string) (tdown func(t *testing.T)) {
	// if the index already exist then delete it.
	client, err := elasticsearch.NewClient()
	require.NoError(t, err)
	isExist, err := client.IndexExists(indexName).Do(context.Background())
	require.NoError(t, err)

	if isExist {
		// Delete the index
		teardown(t, indexName)
	}
	return func(t *testing.T) {
		t.Log("test: tearing down")
		teardown(t, indexName)
	}
}

func teardown(t *testing.T, indexName string) {
	client, err := elasticsearch.NewClient()
	require.NoError(t, err)
	_, err = client.DeleteIndex(indexName).Do(context.Background())
	require.NoError(t, err)
}
