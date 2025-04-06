// Package wrapper used for copy http.ResponseWriter access.
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

// New function creates a wrapper for the http.ResponseWriter.
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

// Write function overwrites the http.ResponseWriter Write() function.
func (rww *Wrapper) Write(buf []byte) (int, error) {
	rww.r.Body = io.NopCloser(bytes.NewBuffer(buf))

	return rww.w.Write(buf)
}

// Header function overwrites the http.ResponseWriter Header() function.
func (rww *Wrapper) Header() http.Header {
	return rww.w.Header()
}

// WriteHeader function overwrites the http.ResponseWriter WriteHeader() function.
func (rww *Wrapper) WriteHeader(statusCode int) {
	rww.r.StatusCode = statusCode

	rww.w.WriteHeader(statusCode)
}

// Result returns response.
func (rww *Wrapper) Result() *http.Response {
	return &rww.r
}
