// Package wrapper used for copy [http.ResponseWriter] access.
package wrapper

import (
	"bytes"
	"io"
	"net/http"
)

// Wrapper struct is used to log the response.
type Wrapper struct {
	w http.ResponseWriter
	r http.Response
}

// New function creates a wrapper for the [http.ResponseWriter].
func New(w http.ResponseWriter, r *http.Request) Wrapper {
	return Wrapper{
		w: w,
		r: http.Response{
			StatusCode: http.StatusOK,
			Request:    r,
			Header:     w.Header(),
		},
	}
}

// Write function overwrites the [http.ResponseWriter] Write() function.
func (w *Wrapper) Write(buf []byte) (int, error) {
	w.r.Body = io.NopCloser(bytes.NewBuffer(buf))

	return w.w.Write(buf)
}

// Header function overwrites the [http.ResponseWriter] Header() function.
func (w *Wrapper) Header() http.Header {
	return w.w.Header()
}

// WriteHeader function overwrites the [http.ResponseWriter] WriteHeader() function.
func (w *Wrapper) WriteHeader(statusCode int) {
	w.r.StatusCode = statusCode

	w.w.WriteHeader(statusCode)
}

// Result returns response.
func (w *Wrapper) Result() *http.Response {
	return &w.r
}
