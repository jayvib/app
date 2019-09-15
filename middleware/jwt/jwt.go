package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Authenticate(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := request.ParseFromRequest(c.Request, ExtractorFunc(TokenFromCookie), func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		// TODO: Redirect to login page.
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
		}
		c.Next()
	}
}

type ExtractorFunc func(r *http.Request) (string, error)

func (e ExtractorFunc) ExtractToken(r *http.Request) (string, error) {
	return e(r)
}

func TokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
