// Package dumper provides server-side HTTP traffic dumping functionality.
package dumper

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/nafigator/http/headers"
	"github.com/nafigator/http/mime"
	"github.com/nafigator/http/response/wrapper"
)

const (
	defaultTemplate = "HTTP dump:\n%s\n\n%s\n"
)

type masker interface {
	Mask(*http.Request, *string)
}

type flusher interface {
	Flush(ctx context.Context, msg string)
}

type logger interface {
	Error(args ...interface{})
}

type HTTPDumper struct {
	masker   masker
	flusher  flusher
	log      logger
	filter   func(string) bool
	template string
}

// New creates http-dumper instance.
func New(
	flusher flusher,
) *HTTPDumper {
	return &HTTPDumper{
		template: defaultTemplate,
		flusher:  flusher,
		filter:   needBody,
	}
}

// WithTemplate initializes new custom output template.
func (h *HTTPDumper) WithTemplate(t string) *HTTPDumper {
	h.template = t

	return h
}

// WithMasker initializes sensitive data masker for dumper output.
func (h *HTTPDumper) WithMasker(m masker) *HTTPDumper {
	h.masker = m

	return h
}

// WithErrLogger initializes logger for dump errors.
func (h *HTTPDumper) WithErrLogger(log logger) *HTTPDumper {
	h.log = log

	return h
}

// WithFilter replaces MIME-based filter function to custom one.
func (h *HTTPDumper) WithFilter(f func(string) bool) *HTTPDumper {
	h.filter = f

	return h
}

func (h *HTTPDumper) MiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, e := httputil.DumpRequest(r, h.filter(r.Header.Get(headers.ContentType)))
		if e != nil {
			if h.log != nil {
				h.log.Error("HTTP request dump error: ", e)
			}

			dump := ""

			h.handleRequest(w, r, next, dump)

			return
		}

		dump := string(b)
		if h.masker != nil {
			h.masker.Mask(r, &dump)
		}

		h.handleRequest(w, r, next, dump)
	})
}

func (h *HTTPDumper) handleRequest(w http.ResponseWriter, r *http.Request, next http.Handler, reqDump string) {
	ctx := r.Context()
	ww := wrapper.New(w, r)

	// Process request
	next.ServeHTTP(&ww, r)

	res := ww.Result()
	res.ProtoMinor = r.ProtoMinor
	res.ProtoMajor = r.ProtoMajor
	defer func() {
		if res.Body != nil {
			_ = res.Body.Close()
		}
	}()

	b, e := httputil.DumpResponse(res, h.filter(res.Header.Get(headers.ContentType)))
	if e != nil {
		if h.log != nil {
			h.log.Error("HTTP response dump error: ", e)
		}

		msg := fmt.Sprintf(h.template, reqDump, "")
		h.flusher.Flush(ctx, msg)

		return
	}

	resDump := string(b)
	if h.masker != nil {
		h.masker.Mask(r, &resDump)
	}

	msg := fmt.Sprintf(h.template, reqDump, resDump)
	h.flusher.Flush(ctx, msg)
}

func needBody(ct string) bool {
	if ct == mime.Bin || ct == "" {
		return false // do not dump files
	}

	return true
}
