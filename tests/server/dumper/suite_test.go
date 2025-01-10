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
	patchDumpRequest patch = iota + 1
	patchDumpResponse
)

type patch uint8

type suite struct {
	ss.Suite
}

// TestRun run tests suite.
func TestRun(t *testing.T) {
	ss.Run(t, &suite{})
}

func applyPatch(p patch) *monkey.PatchGuard {
	switch p {
	case patchDumpResponse:
		return monkey.Patch(httputil.DumpResponse, func(_ *http.Response, _ bool) ([]byte, error) {
			return nil, errors.New("dump response error")
		})
	case patchDumpRequest:
		return monkey.Patch(httputil.DumpRequest, func(_ *http.Request, _ bool) ([]byte, error) {
			return nil, errors.New("dump request error")
		})
	}

	return nil
}
