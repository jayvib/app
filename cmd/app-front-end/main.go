package main

import (
	"github.com/gin-gonic/gin"
	userhttp "github.com/jayvib/app/user/delivery/http"
	"log"
	"net/http"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	router := gin.Default()
	router.StaticFS("/assets", gin.Dir("web/app/assets", false))
	router.LoadHTMLGlob("web/app/templates/**/*.tmpl")
	router.GET("/", func(c *gin.Context){
		c.HTML(http.StatusOK, "index", gin.H{
			"title": "App",
			"message": "Welcome to App!",
		})
	})

	userhttp.RegisterViewHandlers(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
