package json_test

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/nafigator/http/masker/auth"
	"github.com/nafigator/http/masker/json"
	"github.com/stretchr/testify/assert"
)

type next interface {
	Mask(*http.Request, *string)
}

type testCase struct {
	request  http.Request
	next     next
	dump     string
	params   []string
	unmasked *int
	expected string
}

func TestMask(t *testing.T) {
	for name, c := range dataProvider() {
		t.Run(name, func(t *testing.T) {
			s := json.New(c.params).WithNext(c.next)

			if c.unmasked != nil {
				s.WithUnmasked(*c.unmasked)
			}

			s.Mask(&c.request, &c.dump)

			assert.Equal(t, c.expected, c.dump, "Unexpected mask result")
		})
	}
}

func dataProvider() map[string]testCase {
	return map[string]testCase{
		"request with string value replacement": {
			request: http.Request{
				ProtoMajor: 1,
				ProtoMinor: 1,
				Method:     http.MethodPost,
				URL: &url.URL{
					Scheme: "http",
					Host:   "avito.ru",
					Path:   "/user/121",
				},
				Header: map[string][]string{
					"Host":         {"avito.ru"},
					"Content-Type": {"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(`{"name":"Николай","password":"mega-super-pass"}`)),
				Host: "avito.ru",
			},
			params:   []string{"password"},
			dump:     "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"name\":\"Николай\",\"password\":\"mega-super-pass\"}\r\n", //nolint:lll
			expected: "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"name\":\"Николай\",\"password\":\"********er-pass\"}\r\n", //nolint:lll
		},
		"request with null value replacement": {
			request: http.Request{
				ProtoMajor: 1,
				ProtoMinor: 1,
				Method:     http.MethodPost,
				URL: &url.URL{
					Scheme: "http",
					Host:   "avito.ru",
					Path:   "/user/121",
				},
				Header: map[string][]string{
					"Host":         {"avito.ru"},
					"Content-Type": {"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(`{"secret-nullable":null}`)),
				Host: "avito.ru",
			},
			unmasked: toPtr(0),
			params:   []string{"secret-nullable"},
			dump:     "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"secret-nullable\":null}\r\n", //nolint:lll
			expected: "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"secret-nullable\":****}\r\n", //nolint:lll
		},
		"request with bool value replacement": {
			request: http.Request{
				ProtoMajor: 1,
				ProtoMinor: 1,
				Method:     http.MethodPost,
				URL: &url.URL{
					Scheme: "http",
					Host:   "avito.ru",
					Path:   "/user/121",
				},
				Header: map[string][]string{
					"Host":         {"avito.ru"},
					"Content-Type": {"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(`{"bool1":true,"bool2":false}`)),
				Host: "avito.ru",
			},
			unmasked: toPtr(0),
			params:   []string{"bool1", "bool2"},
			dump:     "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"bool1\":true,\"bool2\":false}\r\n", //nolint:lll
			expected: "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"bool1\":****,\"bool2\":*****}\r\n", //nolint:lll
		},
		"request with int value replacement": {
			request: http.Request{
				ProtoMajor: 1,
				ProtoMinor: 1,
				Method:     http.MethodPost,
				URL: &url.URL{
					Scheme: "http",
					Host:   "avito.ru",
					Path:   "/user/121",
				},
				Header: map[string][]string{
					"Host":         {"avito.ru"},
					"Content-Type": {"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(`{"secret-int":123456789}`)),
				Host: "avito.ru",
			},
			unmasked: toPtr(2),
			params:   []string{"secret-int"},
			dump:     "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"secret-int\":123456789}\r\n", //nolint:lll
			expected: "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"secret-int\":*******89}\r\n", //nolint:lll
		},
		"request with int value replacement and bearer": {
			request: http.Request{
				ProtoMajor: 1,
				ProtoMinor: 1,
				Method:     http.MethodPost,
				URL: &url.URL{
					Scheme: "http",
					Host:   "avito.ru",
					Path:   "/user/121",
				},
				Header: map[string][]string{
					"Host":          {"avito.ru"},
					"Content-Type":  {"application/json"},
					"Authorization": {"Bearer super-secret-mega-token-forever"},
				},
				Body: io.NopCloser(strings.NewReader(`{"secret-int":123456789}`)),
				Host: "avito.ru",
			},
			unmasked: toPtr(2),
			next:     auth.New(),
			params:   []string{"secret-int"},
			dump:     "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer super-secret-mega-token-forever\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"secret-int\":123456789}\r\n", //nolint:lll
			expected: "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer ************************forever\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"secret-int\":*******89}\r\n", //nolint:lll
		},
		"request with int and high unmasked": {
			request: http.Request{
				ProtoMajor: 1,
				ProtoMinor: 1,
				Method:     http.MethodPost,
				URL: &url.URL{
					Scheme: "http",
					Host:   "avito.ru",
					Path:   "/user/121",
				},
				Header: map[string][]string{
					"Host":         {"avito.ru"},
					"Content-Type": {"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(`{"secret-int":123456789}`)),
				Host: "avito.ru",
			},
			unmasked: toPtr(12345),
			params:   []string{"secret-int"},
			dump:     "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"secret-int\":123456789}\r\n", //nolint:lll
			expected: "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"secret-int\":123456789}\r\n", //nolint:lll
		},
		"request without matches": {
			request: http.Request{
				ProtoMajor: 1,
				ProtoMinor: 1,
				Method:     http.MethodPost,
				URL: &url.URL{
					Scheme: "http",
					Host:   "avito.ru",
					Path:   "/user/121",
				},
				Header: map[string][]string{
					"Host":         {"avito.ru"},
					"Content-Type": {"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(`{"secret-int":123456789}`)),
				Host: "avito.ru",
			},
			unmasked: toPtr(12345),
			params:   []string{"name", "password"},
			dump:     "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"secret-int\":123456789}\r\n", //nolint:lll
			expected: "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"secret-int\":123456789}\r\n", //nolint:lll
		},
	}
}

// toPtr returns pointer to type.
func toPtr[T any](s T) *T {
	return &s
}
