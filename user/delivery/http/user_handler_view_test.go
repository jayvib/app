// +build unit

package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jayvib/app/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoginPage(t *testing.T) {
	e := gin.Default()
	e.LoadHTMLGlob("../../../web/app/templates/**/*.tmpl")
	e.GET("/user/login", LoginPage)
	w := performRequest(e, http.MethodGet, "/user/login", nil)
	if !strings.Contains(w.Body.String(), "Username") {
		t.Error("expecting word 'Username' but nothing found")
	}
}

func TestRegisterNewUserPage(t *testing.T) {
	e := gin.Default()
	e.LoadHTMLGlob("../../../web/app/templates/**/*.tmpl")
	e.GET("/user/new", RegisterNewUserPage)
	w := performRequest(e, http.MethodGet, "/user/new", nil)
	if !strings.Contains(w.Body.String(), "Username") {
		t.Error("expecting word 'Username' but nothing found")
	}
}

func TestRegisterNewUserRequest(t *testing.T) {
	want := model.User{
		ID:       "uniqueid",
		Email:    "luffy.monkey@onepeice.com",
		Username: "luffy.monkey",
		Token:    "1234567890",
	}

	server := testHTTPServer(t, want)
	defer server.Close()

	e := gin.Default()


	e.POST("/user/new", RegisterNewUserRequest(server.Client()))
	payload := `
{
	"username": "luffy.monkey",
	"email": "luffy.monkey@onepiece.com",
	"password": "pirateking"
}`

	reader := strings.NewReader(payload)
	w := performRequest(e, http.MethodPost, "/user/new", reader)
	assertRegisterNewUserRequest(t, w, want)
}

func assertRegisterNewUserRequest(t *testing.T, response *httptest.ResponseRecorder, want interface{}) {
	t.Helper()
	assert.Equal(t, http.StatusOK, response.Code)
	var got model.User
	err := json.NewDecoder(response.Body).Decode(&got)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func testHTTPServer(t *testing.T, want interface{}) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		t.Log("making a rquest")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(want)
		assert.NoError(t, err)
	}))
	return server
}