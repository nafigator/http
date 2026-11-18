package dumper

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	ss "github.com/stretchr/testify/suite"
)

const (
	patchDumpResponse patch = iota + 1
)

type patch uint8

type suite struct {
	ss.Suite
}

// TestRun run tests suite.
func TestRun(t *testing.T) {
	t.Parallel()
	ss.Run(t, &suite{})
}

func applyPatch(p patch) *gomonkey.Patches {
	if p == patchDumpResponse {
		return gomonkey.ApplyFunc(httputil.DumpResponse, func(_ *http.Response, _ bool) ([]byte, error) {
			return nil, errors.New("dump response error")
		})
	}

	return &gomonkey.Patches{}
}
