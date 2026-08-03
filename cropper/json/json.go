// Package json provides cropper functionality for HTTP dumps.
package json //nolint:revive,nolintlint	// Acknowledged

import (
	"fmt"
	"regexp"
)

type Cropper struct {
	params []string
}

// New creates masker instance.
func New(params []string) *Cropper {
	return &Cropper{
		params: params,
	}
}

// Crop crops out JSON params string values.
func (m *Cropper) Crop(dump *string) {
	for _, p := range m.params {
		re := regexp.MustCompile("(\"" + p + "\"\\s*:\\s*\")([^\"]+)(\")")
		matches := re.FindAllStringSubmatch(*dump, -1)

		if matches == nil {
			continue
		}

		prefix := matches[0][1]
		val := matches[0][2]
		suffix := matches[0][3]

		replacement := fmt.Sprintf("cropped %d bytes", len(val))

		*dump = re.ReplaceAllString(*dump, prefix+replacement+suffix)
	}
}
