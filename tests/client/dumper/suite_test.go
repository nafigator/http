package dumper_test

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"testing"

	"bou.ke/monkey"
	ss "github.com/stretchr/testify/suite"
)

const (
	patchDumpResponse patch = iota + 1
)

type suite struct {
	ss.Suite
}

// TestRun run tests suite.
func TestRun(t *testing.T) {
	ss.Run(t, &suite{})
}

func applyPatch(p patch) *monkey.PatchGuard {
	if p == patchDumpResponse {
		return monkey.Patch(httputil.DumpResponse, func(_ *http.Response, _ bool) ([]byte, error) {
			return nil, errors.New("dump response error")
		})
	}

	return nil
}
