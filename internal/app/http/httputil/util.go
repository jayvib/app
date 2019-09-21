package httputil

import (
	"github.com/gin-gonic/gin"
	"github.com/jayvib/app/apperr"
	"net/http"
)

type RequestError struct {
	Msg        string `json:"msg,omitempty"`
	Parameters map[string]string
}

func HandleSearchError(c *gin.Context, err error) {
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
