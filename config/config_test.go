// +build unit

package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew(t *testing.T) {
	conf, err := New()
	require.NoError(t, err)
	require.NotNil(t, conf)
	assert.Equal(t, "mypassword", conf.Database.Pass)
	assert.Len(t, conf.Elasticsearch.Servers, 1)
}
