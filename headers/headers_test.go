package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	a := assert.New(t)

	a.Equal("Content-Type", Normalize("content-type"))
	a.Equal("Content-Type", Normalize("CONTENT-TYPE"))
	a.Equal("Content-Type", Normalize("cONtENT-tYpE"))
}
