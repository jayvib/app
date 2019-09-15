// +build internal

package config

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_loadConfig(t *testing.T) {

	t.Run("Initialize", func(t *testing.T) {
		conf, err := loadConfig()
		require.NoError(t, err)
		require.NotNil(t, conf)
		assert.Equal(t, "mypassword", conf.Database.Pass)
		assert.Len(t, conf.Elasticsearch.Servers, 1)
	})

	// For local development Setup
	t.Run("Development Environment Config", func(t *testing.T) {
		os.Setenv("APP_ENV", "DEVELOPMENT")
		_, err := loadConfig()
		require.NoError(t, err)
		assert.Equal(t, "devel-config.json", getFileName(viper.ConfigFileUsed()))
	})

	// For staging setup
	t.Run("Staging Environment Config", func(t *testing.T) {
		os.Setenv("APP_ENV", "STAGING")
		_, err := loadConfig()
		require.NoError(t, err)
		assert.Equal(t, getFileName(viper.ConfigFileUsed()), "staging-config.json")
	})

	t.Run("Production Environment", func(t *testing.T) {
		os.Setenv("APP_ENV", "PRODUCTION")
		_, err := loadConfig()
		require.NoError(t, err)
		assert.Equal(t, "config.json", getFileName(viper.ConfigFileUsed()))
	})
}

func getFileName(dir string) string {
	_, file := filepath.Split(dir)
	return file
}
