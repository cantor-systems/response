package response_test

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.cantor.systems/response"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

func TestWithStatusOptions(t *testing.T) {
	options := &response.Options{
		StatusData: func(w http.ResponseWriter, r *http.Request, status int) interface{} {
			return map[string]interface{}{"s": status}
		},
	}
	testHandler := &testStatusHandler{
		status: http.StatusTeapot,
	}
	handler := options.Handler(testHandler)

	w := httptest.NewRecorder()
	r := newTestRequest()

	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusTeapot, w.Code)
	assert.Equal(t, w.Body.String(), `{"s":418}`+"\n")
	assert.Equal(t, w.Header().Get("Content-Type"), "application/json; charset=utf-8")
}

func TestWithError(t *testing.T) {
	w := httptest.NewRecorder()
	r := newTestRequest()

	err := errors.New("something went wrong")

	opts := &response.Options{
		Before: func(w http.ResponseWriter, r *http.Request, status int, data interface{}) (int, interface{}) {
			if err, ok := data.(error); ok {
				return status, map[string]interface{}{"error": err.Error()}
			}
			return status, data
		},
	}
	testHandler := &testHandler{
		status: http.StatusInternalServerError,
		data:   err,
	}
	handler := opts.Handler(testHandler)

	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var data map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &data))
	assert.Equal(t, data, map[string]interface{}{"error": err.Error()})
	assert.Equal(t, w.Header().Get("Content-Type"), "application/json; charset=utf-8")
}

func TestBefore(t *testing.T) {
	options := &response.Options{}
	var beforecall map[string]interface{}
	options.Before = func(w http.ResponseWriter, r *http.Request, status int, data interface{}) (int, interface{}) {
		beforecall = map[string]interface{}{
			"w": w, "r": r, "status": status, "data": data,
		}
		return status, data
	}

	testHandler := &testHandler{
		status: http.StatusOK, data: testdata,
	}
	handler := options.Handler(testHandler)

	w := httptest.NewRecorder()
	r := newTestRequest()

	handler.ServeHTTP(w, r)

	assert.Equal(t, beforecall["w"], w)
	assert.Equal(t, beforecall["r"], r)
	assert.Equal(t, beforecall["status"], testHandler.status)
	assert.Equal(t, beforecall["data"], testHandler.data)

	assert.Equal(t, http.StatusOK, w.Code)
	var data map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &data))
	assert.Equal(t, data, testdata)
	assert.Equal(t, w.Header().Get("Content-Type"), "application/json; charset=utf-8")

}

func TestAfter(t *testing.T) {
	options := &response.Options{}
	var aftercall map[string]interface{}
	options.After = func(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
		aftercall = map[string]interface{}{
			"w": w, "r": r, "status": status, "data": data,
		}
	}

	testHandler := &testHandler{
		status: http.StatusOK, data: testdata,
	}
	handler := options.Handler(testHandler)

	w := httptest.NewRecorder()
	r := newTestRequest()

	handler.ServeHTTP(w, r)

	assert.Equal(t, aftercall["w"], w)
	assert.Equal(t, aftercall["r"], r)
	assert.Equal(t, aftercall["status"], testHandler.status)
	assert.Equal(t, aftercall["data"], testHandler.data)

	assert.Equal(t, http.StatusOK, w.Code)
	var data map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &data))
	assert.Equal(t, data, testdata)
	assert.Equal(t, w.Header().Get("Content-Type"), "application/json; charset=utf-8")
}

type testEncoder struct {
	err error
}

func (e *testEncoder) Encode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	io.WriteString(w, "testEncoder")
	return e.err
}

func (e *testEncoder) ContentType(w http.ResponseWriter, r *http.Request) string {
	return "test/encoder"
}

func TestEncoder(t *testing.T) {
	options := &response.Options{}
	var encodercall map[string]interface{}
	options.Encoder = func(w http.ResponseWriter, r *http.Request) response.Encoder {
		encodercall = map[string]interface{}{
			"w": w, "r": r,
		}
		return &testEncoder{}
	}

	testHandler := &testHandler{
		status: http.StatusOK, data: testdata,
	}
	handler := options.Handler(testHandler)

	w := httptest.NewRecorder()
	r := newTestRequest()

	handler.ServeHTTP(w, r)

	assert.Equal(t, encodercall["w"], w)
	assert.Equal(t, encodercall["r"], r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), "testEncoder")
	assert.Equal(t, w.Header().Get("Content-Type"), "test/encoder")

}

func TestEncoderOnErr(t *testing.T) {
	var onErrCall map[string]interface{}
	options := &response.Options{
		OnErr: func(err error) {
			onErrCall = map[string]interface{}{"err": err}
		},
	}
	encoderErr := errors.New("something went wrong while encoding")
	options.Encoder = func(w http.ResponseWriter, r *http.Request) response.Encoder {
		return &testEncoder{
			err: encoderErr,
		}
	}

	testHandler := &testHandler{
		status: http.StatusOK, data: testdata,
	}
	handler := options.Handler(testHandler)

	w := httptest.NewRecorder()
	r := newTestRequest()

	handler.ServeHTTP(w, r)
	assert.Equal(t, onErrCall["err"], encoderErr)

}

func TestMultipleWith(t *testing.T) {
	options := &response.Options{}
	handler := options.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.With(w, r, http.StatusInternalServerError, errors.New("borked"))
		response.With(w, r, http.StatusOK, nil)
	}))

	w := httptest.NewRecorder()
	r := newTestRequest()

	assert.PanicsWithValue(t, "response: multiple responses", func() {
		handler.ServeHTTP(w, r)
	})

	options = &response.Options{
		AllowMultiple: true,
	}
	handler = options.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.With(w, r, http.StatusInternalServerError, errors.New("borked"))
		response.With(w, r, http.StatusOK, nil)
	}))

	w = httptest.NewRecorder()
	r = newTestRequest()

	handler.ServeHTTP(w, r)
}
