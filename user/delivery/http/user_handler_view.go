package http

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jayvib/app/model"
	"net/http"
)

// RESP API Wrapping:
// https://medium.com/@marcus.olsson/writing-a-go-client-for-your-restful-api-c193a2f4998c

const url = "http://localhost:8080/"

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

func RegisterNewUserRequest(client *http.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()
		resp, err := client.Post("/user/new", "application/json", c.Request.Body)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusFailedDependency, RequestError{
				Msg: err.Error(),
			})
			return
		}

		if resp.StatusCode != http.StatusOK {
			c.JSON(resp.StatusCode, RequestError{
				Msg: "Failed to request new user",
			})
			return
		}

		var user model.User
		err = json.NewDecoder(resp.Body).Decode(&user)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusOK, user)
	}
}
