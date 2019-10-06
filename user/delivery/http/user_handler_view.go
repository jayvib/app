package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "user/login", gin.H{
		"title": "Login",
	})
}
