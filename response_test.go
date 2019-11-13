package response_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go.cantor.systems/response"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testdata = map[string]interface{}{"test": true}

func newTestRequest() *http.Request {
	r, err := http.NewRequest("GET", "Something", nil)
	if err != nil {
		panic("bad request: " + err.Error())
	}
	return r
}

func TestWith(t *testing.T) {
	w := httptest.NewRecorder()
	r := newTestRequest()

	response.With(w, r, http.StatusOK, testdata)

	assert.Equal(t, http.StatusOK, w.Code)

	var data map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &data))
	assert.Equal(t, data, testdata)
	assert.Equal(t, w.Header().Get("Content-Type"), "application/json; charset=utf-8")
}

func TestWithStatusDefault(t *testing.T) {
	w := httptest.NewRecorder()
	r := newTestRequest()

	response.WithStatus(w, r, http.StatusTeapot)

	assert.Equal(t, http.StatusTeapot, w.Code)
	var data map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &data))
	assert.Equal(t, data["status"], "I'm a teapot")
	assert.EqualValues(t, data["code"], http.StatusTeapot)
	assert.Equal(t, w.Header().Get("Content-Type"), "application/json; charset=utf-8")
}
