package main

import "github.com/gin-gonic/gin"

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Static("/assets", "web/app/assets")
	router.Static("/tutorial/headFirstHTMLCSS/", "web/tutorial/headFirstHTMLCSS")
	router.Static("/tutorial/headFirstLounge", "web/tutorial/headFirstHTMLCSS/headFirstLounge")
	router.Static("/tutorial/starbuzz/", "web/tutorial/headFirstHTMLCSS/starbuzz")
	router.Run(":8080")
}
