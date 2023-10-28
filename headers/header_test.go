package headers_test

import (
	"testing"

	. "github.com/nafigator/http/headers"
	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("Content-Type", Normalize("content-type"))
	assert.Equal("Content-Type", Normalize("CONTENT-TYPE"))
	assert.Equal("Content-Type", Normalize("cONtENT-tYpE"))
}
