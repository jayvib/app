package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"testing"
)

func TestLoginPage(t *testing.T) {
	e := gin.Default()
	e.LoadHTMLGlob("../../../web/app/templates/user.tmpl")
	e.GET("/view/user/login", LoginPage)
	w := performRequest(e, http.MethodGet,"/view/user/login", nil)
	if !strings.Contains(w.Body.String(), "Username") {
		t.Error("expecting word 'Username' but nothing found")
	}
}