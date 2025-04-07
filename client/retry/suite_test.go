package retry

import (
	"testing"

	ss "github.com/stretchr/testify/suite"
)

type suite struct {
	ss.Suite
}

// TestRun run tests suite.
func TestRun(t *testing.T) {
	ss.Run(t, &suite{})
}
