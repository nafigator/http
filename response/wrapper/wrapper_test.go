package wrapper

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResult(t *testing.T) {
	url := "https://example.net/v1/user/1"
	expectedBody := `{"name":"saul", "lastName":"goodman"}`
	expectedCount := 37
	expectedCode := http.StatusBadRequest
	expectedHeader := http.Header{
		"Content-Type": {"text/plain; charset=utf-8"},
	}

	unexpectedCount := "Unexpected bytes count"
	unexpectedCode := "Unexpected status code"
	unexpectedBody := "Unexpected body"
	unexpectedHeaders := "Unexpected headers"

	a := assert.New(t)
	w := httptest.ResponseRecorder{}
	r, _ := http.NewRequest(http.MethodGet, url, nil)

	rw := New(&w, r)

	actualCount, _ := rw.Write([]byte(expectedBody))
	a.Equal(expectedCount, actualCount, unexpectedCount)

	rw.WriteHeader(expectedCode)

	res := rw.Result()
	a.Equal(expectedCode, res.StatusCode, unexpectedCode)

	actualHeader := rw.Header()
	a.Equal(expectedHeader, actualHeader, unexpectedHeaders)

	actualBody, _ := io.ReadAll(res.Body)
	a.Equal(expectedBody, string(actualBody), unexpectedBody)
}
