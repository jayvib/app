package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jayvib/app/apperr"
	"github.com/jayvib/app/log"
	"github.com/jayvib/app/user"
	"net/http"
	"strconv"
)

func RegisterUserSearchEngine(e *gin.Engine, se user.SearchEngine) {
	userSearchEngineHandler := NewUserSearchEngineHandler(se)
	e.GET("/search/user", userSearchEngineHandler.Search)
}

func NewUserSearchEngineHandler(se user.SearchEngine) *UserSearchHandler {
	return &UserSearchHandler{
		searchEngine: se,
	}
}

type UserSearchHandler struct {
	searchEngine user.SearchEngine
	l            log.Logger
}

func (us *UserSearchHandler) Search(c *gin.Context) {
	query := c.Query("q")
	sizeStr := c.DefaultQuery("size", "10")
	fromStr := c.DefaultQuery("from", "0")
	log.Infof("Query: %s", c.Query("q"))
	log.Infof("Size: %s", c.Query("size"))
	var ctx context.Context
	if ctx = c.Request.Context(); ctx == nil {
		ctx = context.Background()
	}

	parseInt := func(s string) (int, error) {
		v, err := strconv.ParseInt(s, 10, 64)
		return int(v), err
	}

	size, err := parseInt(sizeStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, RequestError{Msg: "'size' parameter must be a number"})
	}
	from, err := parseInt(fromStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, RequestError{Msg: "'from' parameter must be a number"})
	}

	input := user.SearchInput{
		Query: query,
		Size:  size,
		From:  from,
	}
	res, err := us.searchEngine.Search(ctx, input)
	if err != nil {
		handleSearchError(c, err)
	}
	log.Infof("Result: %#v", res)
	c.JSON(http.StatusOK, res)
}

func handleSearchError(c *gin.Context, err error) {
	if aerr, ok := err.(apperr.Error); ok {
		switch aerr.Code() {
		case apperr.NoItemFound:
			c.AbortWithStatusJSON(http.StatusNotFound, RequestError{
				Msg:        aerr.Error(),
				Parameters: aerr.ExtraInfo(),
			})
		case apperr.BadParameter:
			c.AbortWithStatusJSON(http.StatusBadRequest, RequestError{
				Msg:        aerr.Error(),
				Parameters: aerr.ExtraInfo(),
			})
		case apperr.ValidationErr:
			c.AbortWithStatusJSON(http.StatusBadRequest, RequestError{
				Msg:        aerr.Error(),
				Parameters: aerr.ExtraInfo(),
			})
		case apperr.EmptyID:
			c.AbortWithStatusJSON(http.StatusBadRequest, RequestError{
				Msg:        aerr.Error(),
				Parameters: aerr.ExtraInfo(),
			})
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, RequestError{
				Msg:        aerr.Error(),
				Parameters: aerr.ExtraInfo(),
			})
		}
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, RequestError{
			Msg:        aerr.Error(),
			Parameters: aerr.ExtraInfo(),
		})
	}
}
