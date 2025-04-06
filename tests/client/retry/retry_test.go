package retry_test

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/nafigator/http/client/retry"
	"github.com/nafigator/http/headers"
	"github.com/nafigator/http/mime"
	"github.com/nafigator/http/tests/client/retry/mocks"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

const (
	unexpectedMsgCount = "Unexpected messages count"
	unexpectedRecCount = "Unexpected requests count"
	unexpectedResults  = "Unexpected log results"
	unexpectedResponse = "Unexpected response"
	unexpectedError    = "Unexpected error"
	URL                = "https://localhost"
	count              = 4
)

type roundTripMultipleCase struct {
	expectedFirstErr  error
	expectedSecondErr error
	expectedThirdErr  error
	expectedFourthErr error
	request           *http.Request
	expectedFirstRes  *http.Response
	expectedSecondRes *http.Response
	expectedThirdRes  *http.Response
	expectedFourthRes *http.Response
	name              string
	expectedCount     uint
}

type roundTripOnceCase struct {
	expectedErr      error
	ctx              context.Context
	request          *http.Request
	expectedRes      *http.Response
	validator        func(r *http.Response, err error) bool
	name             string
	expected         []observer.LoggedEntry
	expectedMsgCount int
	expectedCount    uint
	timeout          time.Duration
	expectedMsgLevel zapcore.Level
}

func (s *suite) TestRoundTripMultiple() {
	for _, c := range roundTripMultipleProvider() {
		s.Run(c.name, func() {
			ctrl := gomock.NewController(s.T())
			next := mocks.NewMockRoundTripper(ctrl)

			r := retry.New(next).
				WithPause(0).
				WithLimit(count)

			first := next.EXPECT().
				RoundTrip(gomock.Any()).
				Return(c.expectedFirstRes, c.expectedFirstErr)
			second := next.EXPECT().
				RoundTrip(gomock.Any()).
				Return(c.expectedSecondRes, c.expectedSecondErr)
			third := next.EXPECT().
				RoundTrip(gomock.Any()).
				Return(c.expectedThirdRes, c.expectedThirdErr)
			fourth := next.EXPECT().
				RoundTrip(gomock.Any()).
				Return(c.expectedFourthRes, c.expectedFourthErr)

			gomock.InOrder(
				first,
				second,
				third,
				fourth,
			)

			actualResponse, err := r.RoundTrip(c.request)

			s.Equal(c.expectedFourthRes, actualResponse, unexpectedResponse)
			s.Equal(c.expectedCount, r.Count(), unexpectedRecCount)

			s.Require().NoError(err, unexpectedError)
		})

		_ = c.expectedFirstRes.Body.Close()
		_ = c.expectedSecondRes.Body.Close()
		_ = c.expectedThirdRes.Body.Close()
		_ = c.expectedFourthRes.Body.Close()
	}
}

func (s *suite) TestRoundTripOnce() {
	for _, c := range roundTripOnceProvider() {
		s.Run(c.name, func() {
			ctrl := gomock.NewController(s.T())
			next := mocks.NewMockRoundTripper(ctrl)

			ob, logs := observer.New(c.expectedMsgLevel)
			logger := zap.New(ob).Sugar()
			r := retry.New(next).
				WithPause(0).
				WithLimit(count).
				WithErrLogger(logger)

			if c.timeout != 0 {
				r.WithTimeout(c.timeout)
			}

			if c.ctx != nil {
				r.WithCancel(c.ctx)
			}

			if c.validator != nil {
				r.WithRespValidator(c.validator)
			}

			next.EXPECT().
				RoundTrip(gomock.Any()).
				Return(c.expectedRes, c.expectedErr)

			actualResponse, err := r.RoundTrip(c.request)

			actual := logs.AllUntimed()

			s.Len(actual, c.expectedMsgCount, unexpectedMsgCount)
			s.Equal(c.expected, actual, unexpectedResults)
			s.Equal(c.expectedRes, actualResponse, unexpectedResponse)
			s.Equal(c.expectedCount, r.Count(), unexpectedRecCount)

			if c.expectedErr == nil {
				s.Require().NoError(err, unexpectedError)
				return
			}

			s.Require().Error(c.expectedErr)
		})

		if c.expectedRes != nil {
			_ = c.expectedRes.Body.Close()
		}
	}
}

func roundTripMultipleProvider() []roundTripMultipleCase {
	reqBody := []byte(`{"name":"Boris", "age": 20}`)
	request, _ := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(reqBody))
	request.Header.Set(headers.ContentType, mime.JSON)

	response200 := httptest.NewRecorder()
	response200.Code = http.StatusOK
	response200.Body = bytes.NewBufferString(`{"status":"OK"}`)

	response500 := httptest.NewRecorder()
	response500.Code = http.StatusInternalServerError

	response502 := httptest.NewRecorder()
	response502.Code = http.StatusBadGateway

	response503 := httptest.NewRecorder()
	response503.Code = http.StatusServiceUnavailable

	response504 := httptest.NewRecorder()
	response504.Code = http.StatusGatewayTimeout

	return []roundTripMultipleCase{
		{
			name:              "200 response after 3 retries",
			request:           request,
			expectedFirstRes:  response502.Result(),
			expectedSecondRes: response503.Result(),
			expectedThirdRes:  response504.Result(),
			expectedFourthRes: response200.Result(),
			expectedCount:     count,
		},
		{
			name:              "502 responses",
			request:           request,
			expectedFirstRes:  response502.Result(),
			expectedSecondRes: response502.Result(),
			expectedThirdRes:  response502.Result(),
			expectedFourthRes: response502.Result(),
			expectedCount:     count,
		},
	}
}

func roundTripOnceProvider() []roundTripOnceCase {
	request, _ := http.NewRequest(http.MethodPost, URL, bytes.NewBufferString(`{"name":"Boris", "age": 20}`))
	request.Header.Set(headers.ContentType, mime.JSON)

	response200 := httptest.NewRecorder()
	response200.Code = http.StatusOK
	response200.Body = bytes.NewBufferString(`{"status":"OK"}`)

	response500 := httptest.NewRecorder()
	response500.Code = http.StatusInternalServerError

	response502 := httptest.NewRecorder()
	response502.Code = http.StatusBadGateway

	response503 := httptest.NewRecorder()
	response503.Code = http.StatusServiceUnavailable

	response504 := httptest.NewRecorder()
	response504.Code = http.StatusGatewayTimeout

	return []roundTripOnceCase{
		{
			name:          "200 response with timeout",
			request:       request,
			expectedRes:   response200.Result(),
			expected:      []observer.LoggedEntry{},
			expectedCount: 1,
			timeout:       time.Second,
		},
		{
			name:          "200 response with context",
			request:       request,
			expectedRes:   response200.Result(),
			expected:      []observer.LoggedEntry{},
			expectedCount: 1,
			ctx:           context.Background(),
		},
		{
			name:          "200 response with timeout and context",
			request:       request,
			expectedRes:   response200.Result(),
			expected:      []observer.LoggedEntry{},
			expectedCount: 1,
			timeout:       time.Second,
			ctx:           context.Background(),
		},
		{
			name:    "nil response with error log",
			request: request,
			expectedErr: &net.DNSError{
				Err:         "not found",
				Name:        "host",
				Server:      "192.168.1.100",
				IsTimeout:   false,
				IsTemporary: false,
				IsNotFound:  false,
			},
			expected: []observer.LoggedEntry{{
				Entry:   zapcore.Entry{Level: zap.ErrorLevel, Message: "lookup host on 192.168.1.100: not found"},
				Context: []zapcore.Field{},
			}},
			expectedCount:    1,
			expectedMsgCount: 1,
			timeout:          time.Second,
			ctx:              context.Background(),
		},
		{
			name:          "200 response with custom validator",
			request:       request,
			expectedRes:   response200.Result(),
			expected:      []observer.LoggedEntry{},
			expectedCount: 1,
			validator: func(_ *http.Response, _ error) bool {
				return true
			},
		},
	}
}
