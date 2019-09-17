// +build unit

package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jayvib/app/log"
	"github.com/jayvib/app/model"
	"github.com/jayvib/app/user"
	"github.com/jayvib/app/user/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSearch(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	searchEngineMock := new(mocks.SearchEngine)
	e := gin.Default()
	result := &user.SearchResult{
		Users: []*model.User{
			{
				ID:        "uniqueid",
				Firstname: "Luffy",
				Lastname:  "Monkey",
				Email:     "luffy.monkey@gmail.com",
				Username:  "luff.monkey",
				Password:  "pirateking",
			},
			{
				ID:        "uniquei2",
				Firstname: "Sanji",
				Lastname:  "Vinsmoke",
				Email:     "sanji.vinsmoke@gmail.com",
				Username:  "sanji.vinsmoke",
				Password:  "greatcook",
			},
		},
		TotalHits: 2,
		Next:      0,
	}

	RegisterUserSearchEngine(e, searchEngineMock)
	searchEngineMock.On("Search", mock.Anything, mock.AnythingOfType("user.SearchInput")).
		Return(result, nil).Once()
	w := performRequest(e, http.MethodGet, "/search/user?q=firstname=luffy&size=10", nil)
	if !assert.Equal(t, http.StatusOK, w.Code) {
		printResponseBody(t, w.Body)
	}
	var res user.SearchResult
	err := json.NewDecoder(w.Body).Decode(&res)
	require.NoError(t, err)
	assert.Len(t, res.Users, 2)
	assert.Equal(t, res.TotalHits, 2)
	assert.Equal(t, res.Next, 0)
	searchEngineMock.AssertExpectations(t)
}

func printResponseBody(t *testing.T, body io.Reader) {
	t.Helper()
	content, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	t.Log(string(content))
}
