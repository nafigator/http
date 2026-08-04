package header

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/nafigator/pointer"
	"github.com/stretchr/testify/assert"

	"github.com/nafigator/http/headers"
)

type testCase struct {
	request  http.Request
	next     next
	header   string
	dump     string
	unmasked *int
	expected string
}

func TestMask(t *testing.T) {
	t.Parallel()

	for name, c := range dataProvider() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			a := New(c.header)

			if c.next != nil {
				a.WithNext(c.next)
			}

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
			header:   headers.Authorization,
			dump:     "API exchange\nGET /user/151 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer super-secret-mega-token-forever\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll	// In test it's ok
			expected: "API exchange\nGET /user/151 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: *******************************forever\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll	// In test it's ok
		},
		"request value length less than default unmasked": {
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
					"Host":         {"avito.ru"},
					"X-Request-Id": {"cd1a-e"},
				},
				Host: "avito.ru",
			},
			header:   headers.Authorization,
			dump:     "API exchange\nGET /user/152 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nX-Request-Id: cd1a-e\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll	// In test it's ok
			expected: "API exchange\nGET /user/152 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nX-Request-Id: cd1a-e\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll	// In test it's ok
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
					"X-Request-Id":  {"11111111-1111-1111-1111-11111111111"},
				},
				Host: "avito.ru",
			},
			header:   headers.Authorization,
			next:     New(headers.XRequestID),
			dump:     "API exchange\nGET /user/153?secret=FA2C1234FFD5&password=mega-superPASS&param=32 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer super-secret-mega-token-forever\r\nX-Request-Id: 11111111-1111-1111-1111-11111111111\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll	// In test it's ok
			expected: "API exchange\nGET /user/153?secret=FA2C1234FFD5&password=mega-superPASS&param=32 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: *******************************forever\r\nX-Request-Id: ****************************1111111\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll	// In test it's ok
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
					"Host":         {"avito.ru"},
					"X-Request-Id": {"11111111-1111-1111-1111-11111111111"},
				},
				Host: "avito.ru",
			},
			header:   headers.XRequestID,
			dump:     "API exchange\nGET /user/154?quote=1&secret=FA2C1234FFD5&password=mega-superPASS&param=32 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nX-Request-Id: 11111111-1111-1111-1111-11111111111\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll	// In test it's ok
			expected: "API exchange\nGET /user/154?quote=1&secret=FA2C1234FFD5&password=mega-superPASS&param=32 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nX-Request-Id: ****************************1111111\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll	// In test it's ok
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
			header:   headers.Authorization,
			unmasked: pointer.New(4),
			dump:     "API exchange\nGET /user/155 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer super-secret-mega-token-forever\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll	// In test it's ok
			expected: "API exchange\nGET /user/155 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: **********************************ever\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll	// In test it's ok
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
			header:   headers.Authorization,
			unmasked: pointer.New(0),
			dump:     "API exchange\nGET /user/155 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer super-secret-mega-token-forever\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll	// In test it's ok
			expected: "API exchange\nGET /user/155 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: **************************************\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll	// In test it's ok
		},
		"request with empty header": {
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
			header:   "",
			dump:     "API exchange\nGET /user/155 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer super-secret-mega-token-forever\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll	// In test it's ok
			expected: "API exchange\nGET /user/155 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer super-secret-mega-token-forever\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll	// In test it's ok
		},
		"request with empty header & next": {
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
					"X-Request-Id":  {"11111111-1111-1111-1111-11111111111"},
				},
				Host: "avito.ru",
			},
			header:   "",
			next:     New(headers.XRequestID),
			dump:     "API exchange\nGET /user/155 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer super-secret-mega-token-forever\r\nX-Request-Id: 11111111-1111-1111-1111-11111111111\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll	// In test it's ok
			expected: "API exchange\nGET /user/155 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer super-secret-mega-token-forever\r\nX-Request-Id: ****************************1111111\r\nAccept-Encoding: gzip\r\n\r\n\n", //nolint:lll	// In test it's ok
		},
	}
}
