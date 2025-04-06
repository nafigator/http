package retry

import (
	"context"
	"math"
	"net/http"
	"time"
)

const (
	defaultLimit = 10
	defaultPause = 30 * time.Second
	Forever      = -1 // Negative limit value causes make MaxUint attempts.
)

type logger interface {
	Error(args ...interface{})
}

type HTTPRetry struct {
	ctx      context.Context
	next     http.RoundTripper
	log      logger
	validate func(res *http.Response, err error) bool
	limit    int
	current  uint
	done     uint
	pause    time.Duration
	timeout  time.Duration
}

// New creates http retry repeater instance.
func New(
	next http.RoundTripper,
) *HTTPRetry {
	return &HTTPRetry{
		next:    next,
		limit:   defaultLimit,
		current: 1,
		pause:   defaultPause,
		validate: func(res *http.Response, _ error) bool {
			if res != nil && (res.StatusCode == http.StatusBadGateway ||
				res.StatusCode == http.StatusServiceUnavailable ||
				res.StatusCode == http.StatusGatewayTimeout) {
				return false
			}

			return true
		},
	}
}

// RoundTrip http.RoundTripper implementation.
func (h *HTTPRetry) RoundTrip(req *http.Request) (*http.Response, error) {
	var err error
	var res *http.Response

	for h.checkLimit() {
		if h.current > 1 {
			time.Sleep(h.pause)
		}

		res, err = h.doRequest(req)
		if err != nil && h.log != nil {
			h.log.Error(err.Error())
		}

		h.done++

		if h.validate(res, err) {
			return res, nil
		}

		h.current++
	}

	return res, err
}

// Count returns retries count.
func (h *HTTPRetry) Count() uint {
	return h.done
}

// WithLimit sets retry limit.
func (h *HTTPRetry) WithLimit(limit int) *HTTPRetry {
	h.limit = limit

	return h
}

// WithPause sets retry pause.
func (h *HTTPRetry) WithPause(pause time.Duration) *HTTPRetry {
	h.pause = pause

	return h
}

// WithTimeout sets requests timeouts.
func (h *HTTPRetry) WithTimeout(timeout time.Duration) *HTTPRetry {
	h.timeout = timeout

	return h
}

// WithCancel sets requests contexts.
func (h *HTTPRetry) WithCancel(ctx context.Context) *HTTPRetry {
	h.ctx = ctx

	return h
}

// WithRespValidator sets function for HTTP response validation. Retry on false return.
func (h *HTTPRetry) WithRespValidator(f func(res *http.Response, err error) bool) *HTTPRetry {
	h.validate = f

	return h
}

// WithErrLogger sets logger for response errors.
func (h *HTTPRetry) WithErrLogger(log logger) *HTTPRetry {
	h.log = log

	return h
}

func (h *HTTPRetry) checkLimit() bool {
	if h.current <= uint(h.limit) || (h.limit < 0 && h.current <= math.MaxUint) {
		return true
	}

	return false
}

func (h *HTTPRetry) doRequest(req *http.Request) (*http.Response, error) {
	switch {
	case h.ctx == nil && h.timeout == 0:
		req = req.Clone(context.Background())
	case h.ctx == nil && h.timeout != 0:
		ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
		defer cancel()
		req = req.Clone(ctx)
	case h.ctx != nil && h.timeout != 0:
		ctx, cancel := context.WithTimeout(h.ctx, h.timeout)
		defer cancel()
		req = req.Clone(ctx)
	case h.ctx != nil && h.timeout == 0:
		ctx, cancel := context.WithCancel(h.ctx)
		defer cancel()
		req = req.Clone(ctx)
	}

	return h.next.RoundTrip(req)
}
