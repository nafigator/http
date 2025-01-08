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
	unmasked int
	params   []string
	next     next
}

// New creates masker instance.
func New(params []string, next next) *Masker {
	return &Masker{
		params:   params,
		unmasked: defaultUnmaskedLength,
		next:     next,
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

// LeaveUnmasked sets unmasked chars count at the end of secret.
func (m *Masker) LeaveUnmasked(c int) {
	m.unmasked = c
}
