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
	patchDumpRequest patch = iota + 1
	patchDumpResponse
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
	switch p {
	case patchDumpResponse:
		return gomonkey.ApplyFunc(httputil.DumpResponse, func(_ *http.Response, _ bool) ([]byte, error) {
			return nil, errors.New("dump response error")
		})
	case patchDumpRequest:
		return gomonkey.ApplyFunc(httputil.DumpRequest, func(_ *http.Request, _ bool) ([]byte, error) {
			return nil, errors.New("dump request error")
		})
	}

	return &gomonkey.Patches{}
}
