package http

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jayvib/app/apperr"
	"github.com/jayvib/app/article"
	"github.com/jayvib/app/config"
	"github.com/jayvib/app/model"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func RegisterHandlers(conf *config.Config, r gin.IRouter, u article.Usecase) {
	h := NewHandler(conf, u)
	registerHandlers(r, h)
}

func registerHandlers(r gin.IRouter, h *ArticleHandler) {
	r.GET("/articles", h.FetchArticle)
	r.POST("/articles", h.Store)
	r.GET("/articles/:id", h.GetByID)
	r.DELETE("/articles/:id", h.Delete)
	r.POST("/articles/update", h.Update)
}

type RequestError struct {
	Msg string `json:"msg,omitempty"`
}

func NewHandler(conf *config.Config, uc article.Usecase) *ArticleHandler {
	return &ArticleHandler{
		AUsecase: uc,
		config:   conf,
	}
}

type ArticleHandler struct {
	AUsecase article.Usecase
	config   *config.Config
}

func (a *ArticleHandler) FetchArticle(c *gin.Context) {
	numStr := c.Query("num")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, RequestError{Msg: "num param must be a digit"})
		return
	}
	cursor := c.Query("cursor")
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	listAr, nextCursor, err := a.AUsecase.Fetch(ctx, cursor, num)
	if err != nil {
		c.JSON(getStatusCode(err), RequestError{Msg: err.Error()})
		return
	}
	logrus.Debug("Next cursor:", nextCursor)
	c.Writer.Header().Set("X-Cursor", nextCursor)
	c.JSON(http.StatusOK, listAr)
}

func (a *ArticleHandler) GetByID(c *gin.Context) {
	_id := c.Param("id")
	logrus.Debug("ID:", _id)
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	res, err := a.AUsecase.GetByID(ctx, _id)
	if err != nil {
		c.JSON(getStatusCode(err), RequestError{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (a *ArticleHandler) Store(c *gin.Context) {
	var _article model.Article
	err := json.NewDecoder(c.Request.Body).Decode(&_article)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	logrus.WithFields(logrus.Fields{
		"article": _article,
	}).Debug("Received article")

	if err := a.AUsecase.Store(ctx, &_article); err != nil {
		c.JSON(getStatusCode(err), err.Error())
		return
	}
	c.JSON(http.StatusCreated, _article)
}

func (a *ArticleHandler) Delete(c *gin.Context) {
	_id := c.Param("id")
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	if err := a.AUsecase.Delete(ctx, _id); err != nil {
		c.JSON(getStatusCode(err), RequestError{Msg: err.Error()})
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}

func (a *ArticleHandler) Update(c *gin.Context) {
	var _article model.Article
	if err := json.NewDecoder(c.Request.Body).Decode(&_article); err != nil {
		c.JSON(http.StatusBadRequest, RequestError{Msg: err.Error()})
		return
	}

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	if err := a.AUsecase.Update(ctx, &_article); err != nil {
		c.JSON(getStatusCode(err), RequestError{Msg: err.Error()})
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

func getStatusCode(err error) int {

	if err == nil {
		return http.StatusOK
	}

	switch err {
	case apperr.ItemNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}

}
