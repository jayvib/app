package http

import (
	"github.com/gin-gonic/gin"
	"github.com/jayvib/app/article"
	"github.com/jayvib/app/internal/app/http/httputil"
	"github.com/jayvib/app/internal/app/search"
	"github.com/jayvib/app/log"
	"net/http"
	"strconv"
)

// Please Implement Me
func RegisterSearchHandler(r gin.IRouter, se article.SearchEngine) {
	articleHandler := NewSearchHandler(se)
	r.GET("/search/article", articleHandler.Search)
}

// NewSearchHandler takes a SearchEngine and return the object containing
// all the handlers of the article search.
func NewSearchHandler(se article.SearchEngine) *ArticleSearchHandler {
	return &ArticleSearchHandler{
		se: se,
	}
}

// ArticleSearchHandler contains all the handlers for all
// the search operation.
type ArticleSearchHandler struct {
	se article.SearchEngine
}

// Search is a handler that performs uri query.
func (a *ArticleSearchHandler) Search(c *gin.Context) {
	query := c.Query("q")
	sizeStr := c.DefaultQuery("size", "10")
	fromStr := c.DefaultQuery("from", "0")
	log.Infof("Query: %s Size: %s From: %s", query, sizeStr, fromStr)

	parseInt := func(s string) int {
		v, _ := strconv.ParseInt(s, 10, 64)
		return int(v)
	}

	input := search.Input{
		Query: query,
		Size:  parseInt(sizeStr),
		From:  parseInt(fromStr),
	}
	ctx := c.Request.Context()
	result, err := a.se.Search(ctx, input)
	if err != nil {
		httputil.HandleSearchError(c, err)
	}
	c.JSON(http.StatusOK, result)
}
