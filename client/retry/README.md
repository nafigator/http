<a id="readme-top"></a>
# Go HTTP-client retry

[![GitHub release][Release img]][Release src] [![Github main status][Github main status badge]][Github main status src] [![Go Report Card][Go Report Card badge]][Go Report Card src] [![Coverage report][Codecov report badge]][Codecov report src]

## Features
* Easy of use
* Rich features
* Customizable

## Usage

```go
import (
	"http"
	
    "github.com/nafigator/http/client/retry"
)

...
    r := retry.New(http.DefaultTransport).
        WithPause(RequestPause).
        WithLimit(RequestCount).
        WithTimeout(HTTPTimeout).
        WithErrLogger(log).
        WithCancel(ctx).
        WithRespValidator(func(req *http.Response, _ error) bool {
            return req != nil && req.StatusCode == http.StatusOK
        })
    
    c := http.Client{Transport: r}
    
    resp, err = c.Get("https://example.io/api/v3/checks/")
    if err != nil {
        log.Error(err)
    
        return
    }
...
```
Example above will do `RequestCount` request attempts until 200 response received with `RequestPause` between them.

## Tests
Clone repo and run:
```shell
go test -C tests ./...
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>

[Release img]: https://img.shields.io/github/v/tag/nafigator/http?logo=github&labelColor=333&color=teal&filter=client/retry*
[Release src]: https://github.com/nafigator/http/tree/main/client/retry
[Github main status src]: https://github.com/nafigator/http/tree/main/client/retry
[Github main status badge]: https://github.com/nafigator/http/actions/workflows/go.yml/badge.svg?branch=main
[Go Report Card src]: https://goreportcard.com/report/github.com/nafigator/http/client/retry
[Go Report Card badge]: https://goreportcard.com/badge/github.com/nafigator/http/client/retry
[Codecov report src]: https://app.codecov.io/gh/nafigator/http/tree/main
[Codecov report badge]: https://codecov.io/gh/nafigator/http/branch/main/graph/badge.svg
