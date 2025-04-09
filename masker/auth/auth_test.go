package auth

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	request  http.Request
	next     next
	dump     string
	unmasked *int
	expected string
}

func TestMask(t *testing.T) {
	for name, c := range dataProvider() {
		t.Run(name, func(t *testing.T) {
			a := New().WithNext(c.next)

			if c.unmasked != nil {
				a.WithUnmasked(*c.unmasked)
			}

			a.Mask(&c.request, &c.dump)

			assert.Equal(t, c.expected, c.dump, "Unexpected mask result")
		})
	}
}

func dataProvider() map[string]testCase {
	return map[string]testCase{
		"request with bearer": {
			request: http.Request{
				ProtoMajor: 1,
				ProtoMinor: 1,
				Method:     http.MethodGet,
				URL: &url.URL{
					Scheme: "http",
					Host:   "avito.ru",
					Path:   "/user/151",
				},
				Header: map[string][]string{
					"Host":          {"avito.ru"},
					"Authorization": {"Bearer super-secret-mega-token-forever"},
				},
				Host: "avito.ru",
			},
			dump:     "API exchange\nGET /user/151 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer super-secret-mega-token-forever\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
			expected: "API exchange\nGET /user/151 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer ************************forever\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
		},
		"request with bearer and zero replacement length": {
			request: http.Request{
				ProtoMajor: 1,
				ProtoMinor: 1,
				Method:     http.MethodGet,
				URL: &url.URL{
					Scheme: "http",
					Host:   "avito.ru",
					Path:   "/user/152",
				},
				Header: map[string][]string{
					"Host":          {"avito.ru"},
					"Authorization": {"Bearer token"},
				},
				Host: "avito.ru",
			},
			dump:     "API exchange\nGET /user/152 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer token\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
			expected: "API exchange\nGET /user/152 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer token\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
		},
		"request with bearer and query": {
			request: http.Request{
				ProtoMajor: 1,
				ProtoMinor: 1,
				Method:     http.MethodGet,
				URL: &url.URL{
					Scheme:   "http",
					Host:     "avito.ru",
					Path:     "/user/153",
					RawQuery: "secret=FA2C1234FFD5&password=mega-superPASS&param=32",
				},
				Header: map[string][]string{
					"Host":          {"avito.ru"},
					"Authorization": {"Bearer super-secret-mega-token-forever"},
				},
				Host: "avito.ru",
			},
			next:     New(),
			dump:     "API exchange\nGET /user/153?secret=FA2C1234FFD5&password=mega-superPASS&param=32 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer super-secret-mega-token-forever\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
			expected: "API exchange\nGET /user/153?secret=FA2C1234FFD5&password=mega-superPASS&param=32 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer ************************forever\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
		},
		"request with query and without bearer": {
			request: http.Request{
				ProtoMajor: 1,
				ProtoMinor: 1,
				Method:     http.MethodGet,
				URL: &url.URL{
					Scheme:   "http",
					Host:     "avito.ru",
					Path:     "/user/154",
					RawQuery: "quote=1&secret=FA2C1234FFD5&password=mega-superPASS&param=32",
				},
				Header: map[string][]string{
					"Host": {"avito.ru"},
				},
				Host: "avito.ru",
			},
			next:     New(),
			dump:     "API exchange\nGET /user/154?quote=1&secret=FA2C1234FFD5&password=mega-superPASS&param=32 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
			expected: "API exchange\nGET /user/154?quote=1&secret=FA2C1234FFD5&password=mega-superPASS&param=32 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
		},
		"request with bearer non default unmasked length": {
			request: http.Request{
				ProtoMajor: 1,
				ProtoMinor: 1,
				Method:     http.MethodGet,
				URL: &url.URL{
					Scheme: "http",
					Host:   "avito.ru",
					Path:   "/user/155",
				},
				Header: map[string][]string{
					"Host":          {"avito.ru"},
					"Authorization": {"Bearer super-secret-mega-token-forever"},
				},
				Host: "avito.ru",
			},
			unmasked: toPtr(4),
			dump:     "API exchange\nGET /user/155 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer super-secret-mega-token-forever\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
			expected: "API exchange\nGET /user/155 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer ***************************ever\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
		},
		"request with bearer fully masked": {
			request: http.Request{
				ProtoMajor: 1,
				ProtoMinor: 1,
				Method:     http.MethodGet,
				URL: &url.URL{
					Scheme: "http",
					Host:   "avito.ru",
					Path:   "/user/155",
				},
				Header: map[string][]string{
					"Host":          {"avito.ru"},
					"Authorization": {"Bearer super-secret-mega-token-forever"},
				},
				Host: "avito.ru",
			},
			unmasked: toPtr(0),
			dump:     "API exchange\nGET /user/155 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer super-secret-mega-token-forever\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
			expected: "API exchange\nGET /user/155 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer *******************************\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll
		},
	}
}

// toPtr returns pointer to type.
func toPtr[T any](s T) *T {
	return &s
}
