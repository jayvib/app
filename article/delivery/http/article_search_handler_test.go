// +build unit

package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jayvib/app/article/mocks"
	"github.com/jayvib/app/internal/app/search"
	"github.com/jayvib/app/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestSearch(t *testing.T) {
	r := gin.Default()
	searchEngineMock := new(mocks.SearchEngine)

	expected := &search.Result{
		TotalHits: 2,
		Data: []interface{}{
			&model.Article{
				ID:      "uniqueid1",
				Title:   "Testing 1",
				Content: "This is a content body for testing 1",
			},
			&model.Article{
				ID:      "uniqueid2",
				Title:   "Testing 2",
				Content: "This is a content body for testing 2",
			},
		},
	}

	searchEngineMock.On("Search", mock.Anything, mock.AnythingOfType("search.Input")).
		Return(expected, nil).Once()

	RegisterSearchHandler(r, searchEngineMock)
	w := performRequest(r, http.MethodGet, "/search/article?q=testing", nil)
	require.Equal(t, http.StatusOK, w.Code)

	var got search.Result
	err := json.NewDecoder(w.Body).Decode(&got)
	require.NoError(t, err)
	assert.Equal(t, expected.TotalHits, got.TotalHits)
}
