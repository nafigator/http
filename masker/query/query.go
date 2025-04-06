// Package query provides masking functionality for HTTP dumps.
package query

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
	params   []string
	unmasked int
}

// New creates masker instance.
func New(params []string) *Masker {
	return &Masker{
		params:   params,
		unmasked: defaultUnmaskedLength,
	}
}

// Mask masks value of query-params.
func (m *Masker) Mask(req *http.Request, dump *string) {
	for _, p := range m.params {
		re := regexp.MustCompile(p + "=([^&\\s]+)")
		matches := re.FindAllStringSubmatch(req.URL.RawQuery, -1)

		if matches == nil {
			continue
		}

		val := matches[0][1]

		replacementLength := len(val) - m.unmasked

		if replacementLength < 0 {
			replacementLength = 0
		}

		replacement := strings.Repeat("*", replacementLength) + val[replacementLength:]

		*dump = re.ReplaceAllString(*dump, p+"="+replacement)
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
