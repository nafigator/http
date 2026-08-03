package json //nolint:revive,nolintlint	// Acknowledged

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	dump     string
	expected string
	params   []string
}

func TestCrop(t *testing.T) {
	t.Parallel()

	for name, c := range dataProvider() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s := New(c.params)

			s.Crop(&c.dump)

			assert.Equal(t, c.expected, c.dump, "Unexpected crop result")
		})
	}
}

func dataProvider() map[string]testCase {
	return map[string]testCase{
		"request with string value replacement": {
			params:   []string{"data"},
			dump:     "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"name\":\"Николай\",\"data\":\"data:image/jpeg;base64,/7j/4AAQSkZJRg==\"}\r\n", //nolint:lll	// In test it's ok
			expected: "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"name\":\"Николай\",\"data\":\"cropped 39 bytes\"}\r\n",                        //nolint:lll	// In test it's ok
		},
		"request with null value replacement": {
			params:   []string{"data"},
			dump:     "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"name\":\"Николай\",\"data\":null}\r\n", //nolint:lll	// In test it's ok
			expected: "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"name\":\"Николай\",\"data\":null}\r\n", //nolint:lll	// In test it's ok
		},
		"request with null value": {
			params:   []string{"data"},
			dump:     "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"name\":\"Николай\",\"data\":true}\r\n", //nolint:lll	// In test it's ok
			expected: "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"name\":\"Николай\",\"data\":true}\r\n", //nolint:lll	// In test it's ok
		},
		"request with bool value": {
			params:   []string{"data"},
			dump:     "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"name\":\"Николай\",\"data\":true}\r\n", //nolint:lll	// In test it's ok
			expected: "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"name\":\"Николай\",\"data\":true}\r\n", //nolint:lll	// In test it's ok
		},
		"request with no replacement": {
			params:   []string{"data"},
			dump:     "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"name\":\"Николай\",\"age\":30}\r\n", //nolint:lll	// In test it's ok
			expected: "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"name\":\"Николай\",\"age\":30}\r\n", //nolint:lll	// In test it's ok
		},
		"request with data & date replacement": {
			params:   []string{"data", "date"},
			dump:     "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"name\":\"Николай\",\"data\":\"data:image/jpeg;base64,/7j/4AAQSkZJRg==\",\"date\":\"year:2020, month:July, day:01\"}\r\n", //nolint:lll	// In test it's ok
			expected: "API exchange\nPOST /user/121 HTTP/1.1\r\nHost: avito.ru\r\nUser-Agent: Go-http-client/1.1\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"name\":\"Николай\",\"data\":\"cropped 39 bytes\",\"date\":\"cropped 29 bytes\"}\r\n",                                     //nolint:lll	// In test it's ok
		},
	}
}
