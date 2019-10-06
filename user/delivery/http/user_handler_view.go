package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterViewHandlers(e *gin.Engine) {
	e.GET("/user/login", LoginPage)
	e.GET("/user/new", RegisterNewUserPage)
}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "user/login", gin.H{
		"title": "Login",
	})
}

func RegisterNewUserPage(c *gin.Context) {
	c.HTML(http.StatusOK, "user/new", gin.H{
		"title": "Register",
	})
}
