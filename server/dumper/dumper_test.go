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
	msgOK             = "HTTP dump:\nPOST / HTTP/1.1\r\nHost: localhost\r\nContent-Type: application/json\r\n\r\n{\"name\":\"Boris\", \"age\": 20}\n\nHTTP/1.1 200 OK\r\nContent-Length: 0\r\n\r\n\n"                                         //nolint:lll
	msgOKWithFilter   = "HTTP dump:\nPOST / HTTP/1.1\r\nHost: localhost\r\nContent-Type: application/json\r\n\r\n\n\nHTTP/1.1 200 OK\r\nContent-Length: 0\r\n\r\n\n"                                                                          //nolint:lll
	msgOKWithTemplate = "HTTP dump:\nPOST / HTTP/1.1\r\nHost: localhost\r\nContent-Type: application/json\r\n\r\n{\"name\":\"Boris\", \"age\": 20}\n\n==============\n\nHTTP/1.1 200 OK\r\nContent-Length: 0\r\n\r\n\n"                       //nolint:lll
	internalError     = "HTTP dump:\nPOST / HTTP/1.1\r\nHost: localhost\r\nContent-Type: application/json\r\n\r\n{\"name\":\"Boris\", \"age\": 20}\n\nHTTP/1.1 500 Internal Server Error\r\nContent-Length: 0\r\nConnection: close\r\n\r\n\n" //nolint:lll
	responseDumpError = "HTTP dump:\nPOST / HTTP/1.1\r\nHost: localhost\r\nContent-Type: application/json\r\n\r\n{\"name\":\"Boris\", \"age\": 20}\n\n\n"                                                                                     //nolint:lll
	requestDumpError  = "HTTP dump:\n\n\nHTTP/1.1 200 OK\r\nContent-Length: 0\r\n\r\n\n"
	requestDumpErr    = "HTTP request dump error: dump request error"
	responseDumpErr   = "HTTP response dump error: dump response error"

	unexpectedMsgCount = "Unexpected messages count"
	unexpectedResults  = "Unexpected dump results"
	unexpectedResponse = "Unexpected response"
	unexpectedError    = "Unexpected error"
	URL                = "https://localhost"
)

type handlerCase struct {
	responseRecorder http.ResponseWriter
	masker           masker
	expectedError    error
	request          *http.Request
	filter           func(string) bool
	name             string
	template         string
	expected         []observer.LoggedEntry
	expectedMsgCount int
	usePatch         patch
	expectedMsgLevel zapcore.Level
}

func (s *suite) TestRoundTrip() {
	for _, c := range handlerProvider() {
		s.Run(c.name, func() {
			ctrl := gomock.NewController(s.T())

			ob, logs := observer.New(c.expectedMsgLevel)
			logger := zap.New(ob).Sugar()
			flusher := debug.New(logger)
			d := New(flusher).WithErrLogger(logger)

			if c.template != "" {
				d.WithTemplate(c.template)
			}

			if c.masker != nil {
				d.WithMasker(c.masker)
			}

			if c.filter != nil {
				d.WithFilter(c.filter)
			}

			if c.usePatch > 0 {
				applyPatch(c.usePatch)
				defer func() {
					monkey.UnpatchAll()
				}()
			}

			next := NewMockHandler(ctrl)
			if c.expectedError == nil {
				next.
					EXPECT().
					ServeHTTP(gomock.Any(), c.request).
					Times(1)
			} else {
				next.
					EXPECT().
					ServeHTTP(gomock.Any(), c.request).
					Do(func(next http.ResponseWriter, _ *http.Request) {
						next.Header().Set(headers.Connection, "close")
						next.WriteHeader(http.StatusInternalServerError)
						_, _ = next.Write([]byte("Internal Error"))
					}).
					Times(1)
			}

			d.MiddleWare(next).ServeHTTP(c.responseRecorder, c.request)

			actual := logs.AllUntimed()

			s.Len(actual, c.expectedMsgCount, unexpectedMsgCount)
			s.Equal(c.expected, actual, unexpectedResults)
		})
	}
}

func handlerProvider() []handlerCase {
	reqBody := []byte(`{"name":"Boris", "age": 20}`)
	request, _ := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(reqBody))
	request.Header.Set(headers.ContentType, mime.JSON)

	badRequest := httptest.NewRequest(http.MethodPost, "/ttt", bytes.NewBufferString("Foo"))

	return []handlerCase{
		{
			name:             "200 response",
			request:          request,
			responseRecorder: httptest.NewRecorder(),
			expected: []observer.LoggedEntry{{
				Entry:   zapcore.Entry{Level: zap.DebugLevel, Message: msgOK},
				Context: []zapcore.Field{},
			}},
			expectedMsgLevel: zap.DebugLevel,
			expectedMsgCount: 1,
		},
		{
			name:             "next ServeHTTP writes error",
			request:          request,
			responseRecorder: httptest.NewRecorder(),
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
			usePatch:         patchDumpRequest,
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
			responseRecorder: httptest.NewRecorder(),
			expectedError:    errors.New("internal error"),
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
