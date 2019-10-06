package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("web/app/**/*.tmpl")
	router.GET("/", func(c *gin.Context){
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "App",
			"message": "Welcome to App!",
		})
	})
	router.Run(":8080")
}
