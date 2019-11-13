package response_test

import (
	"go.cantor.systems/response"
	"net/http"
)

var testdata = map[string]interface{}{"test": true}

type testHandler struct {
	status int
	data   interface{}
}

func (t *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response.With(w, r, t.status, t.data)
}

type testStatusHandler struct {
	status int
}

func (t *testStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response.WithStatus(w, r, t.status)
}
