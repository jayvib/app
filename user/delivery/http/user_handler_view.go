package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterViewHandlers(e *gin.Engine) {
	e.GET("/user/login", LoginPage)
}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "user/index", gin.H{
		"title": "Login",
	})
}
