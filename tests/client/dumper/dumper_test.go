package dumper_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"

	"bou.ke/monkey"
	"github.com/nafigator/http/client/dumper"
	"github.com/nafigator/http/headers"
	"github.com/nafigator/http/masker/query"
	"github.com/nafigator/http/storage/debug"
	"github.com/nafigator/http/tests/client/dumper/mocks"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

const (
	msgOK               = "HTTP dump:\nPOST / HTTP/1.1\r\nHost: localhost\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 27\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"name\":\"Boris\", \"age\": 20}\n\nHTTP/1.1 200 OK\r\nConnection: close\r\n\r\n\n" //nolint:lll
	internalError       = "HTTP dump:\nPOST / HTTP/1.1\r\nHost: localhost\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 27\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"name\":\"Boris\", \"age\": 20}\n\ninternal error\n"                               //nolint:lll
	responseDumpError   = "HTTP dump:\nPOST / HTTP/1.1\r\nHost: localhost\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 27\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"name\":\"Boris\", \"age\": 20}\n\n\n"                                             //nolint:lll
	requestDumpError    = "HTTP dump:\n\n\nHTTP/1.1 200 OK\r\nConnection: close\r\n\r\n\n"
	requestDumpWarning  = "HTTP request dump error: unsupported protocol scheme \"\""
	responseDumpWarning = "HTTP response dump error: dump response error"

	unexpectedMsgCount = "Unexpected messages count"
	unexpectedResults  = "Unexpected dump results"
	unexpectedResponse = "Unexpected response"
	unexpectedError    = "Unexpected error"
	URL                = "https://localhost"
	JsonMIME           = "application/json"
)

type masker interface {
	Mask(*http.Request, *string)
}

type patch uint8

type roundRobinCase struct {
	name             string
	request          *http.Request
	responseRecorder *httptest.ResponseRecorder
	usePatch         patch
	masker           masker
	expectedError    error
	expected         []observer.LoggedEntry
	expectedMsgLevel zapcore.Level
	expectedMsgCount int
}

func (s *suite) TestRoundTrip() {
	for _, c := range roundRobinProvider() {
		s.Run(c.name, func() {
			ctrl := gomock.NewController(s.T())
			next := mocks.NewMockRoundTripper(ctrl)

			expectedResponse := c.responseRecorder.Result()

			ob, logs := observer.New(c.expectedMsgLevel)
			logger := zap.New(ob).Sugar()
			flusher := debug.New(logger)
			d := dumper.New("", next, c.masker, flusher, logger)

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

func roundRobinProvider() []roundRobinCase {
	reqBody := []byte(`{"name":"Boris", "age": 20}`)
	request := httptest.NewRequest(http.MethodPost, URL, bytes.NewBuffer(reqBody))
	request.Header.Set(headers.ContentType, JsonMIME)

	errResponse := httptest.NewRecorder()
	errResponse.Code = 500
	errResponse.Body = bytes.NewBuffer([]byte(`{"errors":[{"code":1,"message":"internal error"}]}`))

	badResponse := httptest.NewRecorder()
	badResponse.Code = 500
	badResponse.Body = bytes.NewBuffer([]byte(`{"errors":[{"code":1,"message":"internal error"}]}`))
	_ = badResponse.Result().Body.Close()

	return []roundRobinCase{
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
			request:          httptest.NewRequest(http.MethodPost, "/ttt", bytes.NewBuffer([]byte("Foo"))),
			responseRecorder: httptest.NewRecorder(),
			expectedError:    nil,
			expected: []observer.LoggedEntry{{
				Entry:   zapcore.Entry{Level: zap.WarnLevel, Message: requestDumpWarning},
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
				Entry:   zapcore.Entry{Level: zap.WarnLevel, Message: responseDumpWarning},
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
			masker:           query.New([]string{"password"}, nil),
			expected: []observer.LoggedEntry{{
				Entry:   zapcore.Entry{Level: zap.DebugLevel, Message: msgOK},
				Context: []zapcore.Field{},
			}},
			expectedMsgLevel: zap.DebugLevel,
			expectedMsgCount: 1,
		},
	}
}
