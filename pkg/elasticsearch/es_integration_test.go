// +build integration,elasticsearch

package elasticsearch_test

import (
	"github.com/jayvib/app/config"
	"github.com/jayvib/app/pkg/elasticsearch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := elasticsearch.NewClient()
	require.NoError(t, err)
	require.NotNil(t, client)

	t.Run("Singleton test", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			newClient, err := elasticsearch.NewClient()
			require.NoError(t, err)
			assert.Equal(t, client, newClient)
		}
	})
}

func TestNewSimpleClient(t *testing.T) {
	client, err := elasticsearch.NewSimpleClient()
	require.NoError(t, err)
	require.NotNil(t, client)

	t.Run("Singleton test", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			newClient, err := elasticsearch.NewSimpleClient()
			require.NoError(t, err)
			assert.NotNil(t, newClient)
			assert.Equal(t, client, newClient)
		}
	})
}

func getConfig(t *testing.T) *config.Config {
	conf, err := config.New()
	require.NoError(t, err)
	return conf
}
