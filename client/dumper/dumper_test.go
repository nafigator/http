package dumper

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"

	"bou.ke/monkey"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/nafigator/http/headers"
	"github.com/nafigator/http/masker/query"
	"github.com/nafigator/http/mime"
	"github.com/nafigator/http/storage/debug"
)

const (
	msgOK             = "HTTP dump:\nPOST / HTTP/1.1\r\nHost: localhost\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 27\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"name\":\"Boris\", \"age\": 20}\n\nHTTP/1.1 200 OK\r\nConnection: close\r\n\r\n\n"                   //nolint:lll
	msgOKWithFilter   = "HTTP dump:\nPOST / HTTP/1.1\r\nHost: localhost\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 27\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n\n\nHTTP/1.1 200 OK\r\nConnection: close\r\n\r\n\n"                                                    //nolint:lll
	msgOKWithTemplate = "HTTP dump:\nPOST / HTTP/1.1\r\nHost: localhost\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 27\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"name\":\"Boris\", \"age\": 20}\n\n==============\n\nHTTP/1.1 200 OK\r\nConnection: close\r\n\r\n\n" //nolint:lll
	internalError     = "HTTP dump:\nPOST / HTTP/1.1\r\nHost: localhost\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 27\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"name\":\"Boris\", \"age\": 20}\n\ninternal error\n"                                                 //nolint:lll
	responseDumpError = "HTTP dump:\nPOST / HTTP/1.1\r\nHost: localhost\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 27\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"name\":\"Boris\", \"age\": 20}\n\n\n"                                                               //nolint:lll
	requestDumpError  = "HTTP dump:\n\n\nHTTP/1.1 200 OK\r\nConnection: close\r\n\r\n\n"
	requestDumpErr    = "HTTP request dump error: unsupported protocol scheme \"\""
	responseDumpErr   = "HTTP response dump error: dump response error"

	unexpectedMsgCount = "Unexpected messages count"
	unexpectedResults  = "Unexpected dump results"
	unexpectedResponse = "Unexpected response"
	unexpectedError    = "Unexpected error"
	URL                = "https://localhost"
)

type roundTripCase struct {
	masker           masker
	expectedError    error
	request          *http.Request
	responseRecorder *httptest.ResponseRecorder
	filter           func(string) bool
	name             string
	template         string
	expected         []observer.LoggedEntry
	expectedMsgCount int
	usePatch         patch
	expectedMsgLevel zapcore.Level
}

func (s *suite) TestRoundTrip() {
	for _, c := range roundTripProvider() {
		s.Run(c.name, func() {
			ctrl := gomock.NewController(s.T())
			next := NewMockRoundTripper(ctrl)

			expectedResponse := c.responseRecorder.Result()

			ob, logs := observer.New(c.expectedMsgLevel)
			logger := zap.New(ob).Sugar()
			flusher := debug.New(logger)
			d := New(next, flusher).WithErrLogger(logger)

			if c.template != "" {
				d.WithTemplate(c.template)
			}

			if c.masker != nil {
				d.WithMasker(c.masker)
			}

			if c.filter != nil {
				d.WithFilter(c.filter)
			}

			next.EXPECT().
				RoundTrip(c.request).
				Return(expectedResponse, c.expectedError).
				Times(1)

			if c.usePatch > 0 {
				applyPatch(c.usePatch)
				defer func() {
					monkey.UnpatchAll()
				}()
			}

			actualResponse, err := d.RoundTrip(c.request)

			actual := logs.AllUntimed()

			s.Len(actual, c.expectedMsgCount, unexpectedMsgCount)
			s.Equal(c.expected, actual, unexpectedResults)
			s.Equal(expectedResponse, actualResponse, unexpectedResponse)

			if c.expectedError == nil {
				s.Require().NoError(err, unexpectedError)
				return
			}

			s.Require().Error(c.expectedError)
		})
	}
}

func roundTripProvider() []roundTripCase {
	reqBody := []byte(`{"name":"Boris", "age": 20}`)
	request, _ := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(reqBody))
	request.Header.Set(headers.ContentType, mime.JSON)

	badRequest, _ := http.NewRequest(http.MethodPost, "/ttt", bytes.NewBufferString("Foo"))

	errResponse := httptest.NewRecorder()
	errResponse.Code = 500
	errResponse.Body = bytes.NewBufferString(`{"errors":[{"code":1,"message":"internal error"}]}`)

	badResponse := httptest.NewRecorder()
	badResponse.Code = 500
	badResponse.Body = bytes.NewBufferString(`{"errors":[{"code":1,"message":"internal error"}]}`)
	_ = badResponse.Result().Body.Close()

	return []roundTripCase{
		{
			name:             "200 response",
			request:          request,
			responseRecorder: httptest.NewRecorder(),
			expectedError:    nil,
			expected: []observer.LoggedEntry{{
				Entry:   zapcore.Entry{Level: zap.DebugLevel, Message: msgOK},
				Context: []zapcore.Field{},
			}},
			expectedMsgLevel: zap.DebugLevel,
			expectedMsgCount: 1,
		},
		{
			name:             "next RoundTrip returns error",
			request:          request,
			responseRecorder: errResponse,
			expectedError:    errors.New("internal error"),
			expected: []observer.LoggedEntry{{
				Entry:   zapcore.Entry{Level: zap.DebugLevel, Message: internalError},
				Context: []zapcore.Field{},
			}},
			expectedMsgLevel: zap.DebugLevel,
			expectedMsgCount: 1,
		},
		{
			name:             "dump request error",
			request:          badRequest,
			responseRecorder: httptest.NewRecorder(),
			expectedError:    nil,
			expected: []observer.LoggedEntry{{
				Entry:   zapcore.Entry{Level: zap.ErrorLevel, Message: requestDumpErr},
				Context: []zapcore.Field{},
			}, {
				Entry:   zapcore.Entry{Level: zap.DebugLevel, Message: requestDumpError},
				Context: []zapcore.Field{},
			}},
			expectedMsgLevel: zap.DebugLevel,
			expectedMsgCount: 2,
		},
		{
			name:             "dump response error",
			request:          request,
			responseRecorder: badResponse,
			expectedError:    nil,
			usePatch:         patchDumpResponse,
			expected: []observer.LoggedEntry{{
				Entry:   zapcore.Entry{Level: zap.ErrorLevel, Message: responseDumpErr},
				Context: []zapcore.Field{},
			}, {
				Entry:   zapcore.Entry{Level: zap.DebugLevel, Message: responseDumpError},
				Context: []zapcore.Field{},
			}},
			expectedMsgLevel: zap.DebugLevel,
			expectedMsgCount: 2,
		},
		{
			name:             "200 response with masker",
			request:          request,
			responseRecorder: httptest.NewRecorder(),
			expectedError:    nil,
			masker:           query.New([]string{"password"}),
			expected: []observer.LoggedEntry{{
				Entry:   zapcore.Entry{Level: zap.DebugLevel, Message: msgOK},
				Context: []zapcore.Field{},
			}},
			expectedMsgLevel: zap.DebugLevel,
			expectedMsgCount: 1,
		},
		{
			name:             "200 response with template",
			request:          request,
			responseRecorder: httptest.NewRecorder(),
			expectedError:    nil,
			template:         "HTTP dump:\n%s\n\n==============\n\n%s\n",
			expected: []observer.LoggedEntry{{
				Entry:   zapcore.Entry{Level: zap.DebugLevel, Message: msgOKWithTemplate},
				Context: []zapcore.Field{},
			}},
			expectedMsgLevel: zap.DebugLevel,
			expectedMsgCount: 1,
		},
		{
			name:             "200 response with filter",
			request:          request,
			responseRecorder: httptest.NewRecorder(),
			expectedError:    nil,
			filter: func(ct string) bool {
				return ct != mime.JSON
			},
			expected: []observer.LoggedEntry{{
				Entry:   zapcore.Entry{Level: zap.DebugLevel, Message: msgOKWithFilter},
				Context: []zapcore.Field{},
			}},
			expectedMsgLevel: zap.DebugLevel,
			expectedMsgCount: 1,
		},
	}
}
