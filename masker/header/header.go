// Package header provides header value masking functionality for HTTP dumps.
package header

import (
	"net/http"
	"regexp"
	"strings"
)

const (
	defaultUnmaskedLength = 7
)

type next interface {
	Mask(*http.Request, *string)
}

type Masker struct {
	next     next
	header   string
	unmasked int
}

// New creates masker instance.
func New(h string) *Masker {
	return &Masker{
		unmasked: defaultUnmaskedLength,
		header:   h,
	}
}

// Mask masks value of header.
func (m *Masker) Mask(req *http.Request, dump *string) {
	if m.header == "" {
		if m.next != nil {
			m.next.Mask(req, dump)
		}

		return
	}

	s := req.Header.Get(m.header)
	replacementLength := max(len(s)-m.unmasked, 0)

	s = strings.Repeat("*", replacementLength) + s[replacementLength:]

	re := regexp.MustCompile("(" + regexp.QuoteMeta(m.header) + "\\s*:\\s*)[^\\r]+\\r\\n")
	match := re.FindStringSubmatch(*dump)

	if match != nil {
		*dump = re.ReplaceAllString(*dump, match[1]+s+"\r\n")
	}

	if m.next != nil {
		m.next.Mask(req, dump)
	}
}

// WithNext sets next masker for nested processing.
func (m *Masker) WithNext(n next) *Masker {
	m.next = n

	return m
}

// WithUnmasked sets unmasked chars count at the end of secret.
func (m *Masker) WithUnmasked(c int) *Masker {
	m.unmasked = c

	return m
}
