package dumper

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/nafigator/http/headers"
	"github.com/nafigator/http/mime"
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
	next     http.RoundTripper
	masker   masker
	flusher  flusher
	log      logger
	filter   func(string) bool
	template string
}

// New creates http-dumper instance.
func New(
	next http.RoundTripper,
	flusher flusher,
) *HTTPDumper {
	return &HTTPDumper{
		template: defaultTemplate,
		next:     next,
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

// RoundTrip http.RoundTripper implementation.
func (h *HTTPDumper) RoundTrip(req *http.Request) (*http.Response, error) {
	b, e := httputil.DumpRequestOut(req, h.filter(req.Header.Get(headers.ContentType)))
	if e != nil {
		if h.log != nil {
			h.log.Error("HTTP request dump error: ", e)
		}

		dump := ""

		return h.handleRequest(req, dump)
	}

	dump := string(b)
	if h.masker != nil {
		h.masker.Mask(req, &dump)
	}

	return h.handleRequest(req, dump)
}

func (h *HTTPDumper) handleRequest(req *http.Request, reqDump string) (*http.Response, error) {
	var e error
	var res *http.Response
	ctx := req.Context()

	// Send request
	res, e = h.next.RoundTrip(req)
	if e != nil {
		msg := fmt.Sprintf(h.template, reqDump, e.Error())
		h.flusher.Flush(ctx, msg)

		return res, e
	}

	var b []byte

	b, e = httputil.DumpResponse(res, h.filter(res.Header.Get(headers.ContentType)))
	if e != nil {
		if h.log != nil {
			h.log.Error("HTTP response dump error: ", e)
		}

		msg := fmt.Sprintf(h.template, reqDump, "")
		h.flusher.Flush(ctx, msg)

		return res, nil
	}

	resDump := string(b)
	if h.masker != nil {
		h.masker.Mask(req, &resDump)
	}

	msg := fmt.Sprintf(h.template, reqDump, resDump)
	h.flusher.Flush(ctx, msg)

	return res, e
}

func needBody(ct string) bool {
	if ct == mime.Bin || ct == "" {
		return false // do not dump files
	}

	return true
}
