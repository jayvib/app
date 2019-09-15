package testutil

// testutil provides a convenient functions for testing

import (
	"io"
	"net/http"
	"net/http/httptest"
)

func PerformRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
