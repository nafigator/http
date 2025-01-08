package auth

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/nafigator/http/headers"
)

const (
	defaultUnmaskedLength = 7
	headerTemplate        = "Authorization: "
)

type next interface {
	Mask(*http.Request, *string)
}

type Masker struct {
	unmasked int
	next     next
}

// New creates masker instance.
func New(next next) *Masker {
	return &Masker{
		unmasked: defaultUnmaskedLength,
		next:     next,
	}
}

// Mask masks value of Authorization header.
func (m *Masker) Mask(req *http.Request, dump *string) {
	s := strings.Fields(req.Header.Get(headers.Authorization))
	if len(s) > 0 {
		secretIdx := len(s) - 1
		replacementLength := len(s[secretIdx]) - m.unmasked

		if replacementLength < 0 {
			replacementLength = 0
		}

		s[secretIdx] = strings.Repeat("*", replacementLength) + s[secretIdx][replacementLength:]
	}

	var re = regexp.MustCompile(headerTemplate + ".+\\r\\n")
	*dump = re.ReplaceAllString(*dump, headerTemplate+strings.Join(s, " ")+"\r\n")

	if m.next != nil {
		m.next.Mask(req, dump)
	}
}

// LeaveUnmasked sets unmasked chars count at the end of secret.
func (m *Masker) LeaveUnmasked(c int) {
	m.unmasked = c
}