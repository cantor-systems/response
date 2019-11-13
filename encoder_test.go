package response_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go.cantor.systems/response"
	"net/http/httptest"
	"testing"
)

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	r := newTestRequest()

	assert.Equal(t, response.JSON.ContentType(w, r), "application/json; charset=utf-8")

	assert.NoError(t, response.JSON.Encode(w, r, testdata))

	var data map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &data))
	assert.Equal(t, data, testdata)
}
