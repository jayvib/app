// +build integration,elasticsearch

package http_test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jayvib/app/config"
	estestutil "github.com/jayvib/app/internal/elasticsearch/testutil"
	"github.com/jayvib/app/log"
	"github.com/jayvib/app/model"
	"github.com/jayvib/app/user"
	userhttp "github.com/jayvib/app/user/delivery/http"
	httptestutil "github.com/jayvib/app/user/delivery/http/testutil"
	usersearches "github.com/jayvib/app/user/search/elasticsearch"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

var client *elastic.Client

const (
	elasticsearchTestURL = "http://localhost:9200"
	index                = "user"
)

func init() {
	os.Setenv(config.AppEnvironmentKey, config.DevelopmentEnv)
}

func TestIntegration_Search(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	estestutil.LoadSampleDataFromProvider(t, client, index, provider(t, "users.search.1.input"))

	se := usersearches.New(client)
	e := gin.Default()
	userhttp.RegisterUserSearchEngine(e, se)
	w := httptestutil.PerformRequest(e, http.MethodGet, "/search/user?q=firstname=luffy&size=10", nil)
	require.Equal(t, http.StatusOK, w.Code)

	var res user.SearchResult
	err := json.NewDecoder(w.Body).Decode(&res)
	require.NoError(t, err)
	assert.Len(t, res.Users, 5)
	assert.Equal(t, res.TotalHits, 5)
	assert.Equal(t, res.Next, 0)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	var err error
	client, err = estestutil.Setup(elasticsearchTestURL, index, "user.mapping.es.json")
	if err != nil {
		panic(err)
	}
	exitVal := m.Run()
	os.Exit(exitVal)
}

func provider(t *testing.T, filename string) func() map[string]interface{} {
	t.Helper()
	file, err := os.Open(filepath.Join("testdata", filename))
	require.NoError(t, err)
	defer func() {
		err = file.Close()
		require.NoError(t, err)
	}()

	var users []*model.User
	err = json.NewDecoder(file).Decode(&users)
	require.NoError(t, err)

	userMap := make(map[string]interface{})
	for _, usr := range users {
		userMap[usr.ID] = usr
	}
	return func() map[string]interface{} {
		return userMap
	}
}
