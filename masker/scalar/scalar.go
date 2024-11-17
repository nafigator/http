package scalar

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

// New creates masker instance
func New(params []string, next next) *Masker {
	return &Masker{
		params:   params,
		unmasked: defaultUnmaskedLength,
		next:     next,
	}
}

// Mask masks value of JSON scalars
func (m *Masker) Mask(req *http.Request, dump *string) {
	for _, p := range m.params {
		re := regexp.MustCompile("(\"" + p + "\"\\s*:\\s*\"?)(null|true|false|[\\d]+|[^\"]+)(\")?")
		matches := re.FindAllStringSubmatch(*dump, -1)

		if matches == nil {
			continue
		}

		prefix := matches[0][1]
		val := matches[0][2]
		suffix := matches[0][3]

		replacementLength := len(val) - m.unmasked

		if replacementLength < 0 {
			replacementLength = 0
		}

		replacement := strings.Repeat("*", replacementLength) + val[replacementLength:]

		*dump = re.ReplaceAllString(*dump, prefix+replacement+suffix)
	}

	if m.next != nil {
		m.next.Mask(req, dump)
	}
}

// LeaveUnmasked sets unmasked chars count at the end of secret
func (m *Masker) LeaveUnmasked(c int) {
	m.unmasked = c
}
