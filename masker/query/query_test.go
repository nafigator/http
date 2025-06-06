package query

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	request  http.Request
	next     next
	unmasked *int
	dump     string
	expected string
	params   []string
}

func TestMask(t *testing.T) {
	for name, c := range dataProvider() {
		t.Run(name, func(t *testing.T) {
			q := New(c.params).WithNext(c.next)

			if c.unmasked != nil {
				q.WithUnmasked(*c.unmasked)
			}

			q.Mask(&c.request, &c.dump)

			assert.Equal(t, c.expected, c.dump, "Unexpected mask result")
		})
	}
}

func dataProvider() map[string]testCase {
	return map[string]testCase{
		"request with query": {
			request: http.Request{
				ProtoMajor: 1,
				ProtoMinor: 1,
				Method:     http.MethodGet,
				URL: &url.URL{
					Scheme:   "http",
					Host:     "avito.ru",
					Path:     "/user/151",
					RawQuery: "secret=FA2C1234FFD5&password=mega-superPASS&param=10",
				},
				Header: map[string][]string{
					"Host": {"avito.ru"},
				},
				Host: "avito.ru",
			},
			params:   []string{"password", "secret"},
			dump:     "API exchange\nGET /user/151?secret=FA2C1234FFD5&password=mega-superPASS&param=10 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
			expected: "API exchange\nGET /user/151?secret=*****234FFD5&password=*******perPASS&param=10 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
		},
		"request with query and without matches": {
			request: http.Request{
				ProtoMajor: 1,
				ProtoMinor: 1,
				Method:     http.MethodGet,
				URL: &url.URL{
					Scheme:   "http",
					Host:     "avito.ru",
					Path:     "/user/151",
					RawQuery: "secret=FA2C1234FFD5&password=mega-superPASS&param=10",
				},
				Header: map[string][]string{
					"Host": {"avito.ru"},
				},
				Host: "avito.ru",
			},
			params:   []string{"name", "query"},
			dump:     "API exchange\nGET /user/151?secret=FA2C1234FFD5&password=mega-superPASS&param=10 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
			expected: "API exchange\nGET /user/151?secret=FA2C1234FFD5&password=mega-superPASS&param=10 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
		},
		"request with query and non default unmasked length": {
			request: http.Request{
				ProtoMajor: 1,
				ProtoMinor: 1,
				Method:     http.MethodGet,
				URL: &url.URL{
					Scheme:   "http",
					Host:     "avito.ru",
					Path:     "/user/151",
					RawQuery: "secret=FA2C1234FFD5&password=mega-superPASS&param=10",
				},
				Header: map[string][]string{
					"Host": {"avito.ru"},
				},
				Host: "avito.ru",
			},
			params:   []string{"password", "secret"},
			unmasked: toPtr(3),
			dump:     "API exchange\nGET /user/151?secret=FA2C1234FFD5&password=mega-superPASS&param=10 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
			expected: "API exchange\nGET /user/151?secret=*********FD5&password=***********ASS&param=10 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
		},
		"request with query and zero unmasked length": {
			request: http.Request{
				ProtoMajor: 1,
				ProtoMinor: 1,
				Method:     http.MethodGet,
				URL: &url.URL{
					Scheme:   "http",
					Host:     "avito.ru",
					Path:     "/user/151",
					RawQuery: "secret=FA2C1234FFD5&password=mega-superPASS&param=10",
				},
				Header: map[string][]string{
					"Host": {"avito.ru"},
				},
				Host: "avito.ru",
			},
			params:   []string{"password", "secret"},
			unmasked: toPtr(0),
			dump:     "API exchange\nGET /user/151?secret=FA2C1234FFD5&password=mega-superPASS&param=10 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
			expected: "API exchange\nGET /user/151?secret=************&password=**************&param=10 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
		},
		"request with query and too high unmasked length": {
			request: http.Request{
				ProtoMajor: 1,
				ProtoMinor: 1,
				Method:     http.MethodGet,
				URL: &url.URL{
					Scheme:   "http",
					Host:     "avito.ru",
					Path:     "/user/151",
					RawQuery: "secret=FA2C1234FFD5&password=mega-superPASS&param=10",
				},
				Header: map[string][]string{
					"Host": {"avito.ru"},
				},
				Host: "avito.ru",
			},
			params:   []string{"password", "secret"},
			unmasked: toPtr(1234),
			dump:     "API exchange\nGET /user/151?secret=FA2C1234FFD5&password=mega-superPASS&param=10 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
			expected: "API exchange\nGET /user/151?secret=FA2C1234FFD5&password=mega-superPASS&param=10 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
		},
		"request with query and bearer": {
			request: http.Request{
				ProtoMajor: 1,
				ProtoMinor: 1,
				Method:     http.MethodGet,
				URL: &url.URL{
					Scheme:   "http",
					Host:     "avito.ru",
					Path:     "/user/151",
					RawQuery: "secret=FA2C1234FFD5&password=mega-superPASS&param=10",
				},
				Header: map[string][]string{
					"Host":          {"avito.ru"},
					"Authorization": {"Bearer super-secret-mega-token-forever"},
				},
				Host: "avito.ru",
			},
			next:     New([]string{"type"}),
			params:   []string{"password", "secret"},
			dump:     "API exchange\nGET /user/151?secret=FA2C1234FFD5&password=mega-superPASS&param=10 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer super-secret-mega-token-forever\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
			expected: "API exchange\nGET /user/151?secret=*****234FFD5&password=*******perPASS&param=10 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer super-secret-mega-token-forever\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
		},
	}
}

// toPtr returns pointer to type.
func toPtr[T any](s T) *T {
	return &s
}
