package dumper

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/nafigator/http/headers"
)

const (
	defaultMIME     = "application/octet-stream"
	defaultTemplate = "HTTP dump:\n%s\n\n%s\n"
)

type masker interface {
	Mask(*http.Request, *string)
}

type flusher interface {
	Flush(ctx context.Context, msg string)
}

type logger interface {
	Warn(args ...interface{})
}

type HTTPDumper struct {
	template string
	next     http.RoundTripper
	masker   masker
	flusher  flusher
	log      logger
}

// New creates http-dumper instance.
func New(
	template string,
	next http.RoundTripper,
	masker masker,
	flusher flusher,
	log logger,
) http.RoundTripper {
	if template == "" {
		template = defaultTemplate
	}

	return &HTTPDumper{
		template: template,
		next:     next,
		masker:   masker,
		flusher:  flusher,
		log:      log,
	}
}

// RoundTrip http.RoundTripper implementation.
func (h *HTTPDumper) RoundTrip(req *http.Request) (*http.Response, error) {
	b, e := httputil.DumpRequestOut(req, needBody(req.Header.Get(headers.ContentType)))
	if e != nil {
		h.log.Warn("HTTP request dump error: ", e)

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

	b, e = httputil.DumpResponse(res, needBody(res.Header.Get(headers.ContentType)))
	if e != nil {
		h.log.Warn("HTTP response dump error: ", e)

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
	if ct == defaultMIME || ct == "" {
		return false // do not dump files
	}

	return true
}
