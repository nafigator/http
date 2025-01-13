# response/wrapper

[![GitHub release][Release img]][Release src] [![Github main status][Github main status badge]][Github main status src] [![Go Report Card][Go Report Card badge]][Go Report Card src] [![Coverage report][Codecov report badge]][Codecov report src]

Wrapper for HTTP ResponseWriter.

## Usage

```go
import (
  "fmt"

  "github.com/nafigator/http/response/wrapper"
)

  ...
  rw := wrapper.New(w, r)
  ...
  response := rw.Result()
  ...
```

## Tests
Clone repo and run:
```shell
go test -C tests ./...
```

[Release img]: https://img.shields.io/github/v/tag/nafigator/http?logo=github&labelColor=333&color=teal&filter=response/wrapper*
[Release src]: https://github.com/nafigator/http/tree/main/response/wrapper
[Github main status src]: https://github.com/nafigator/http/tree/main/reponse/wrapper
[Github main status badge]: https://github.com/nafigator/http/actions/workflows/go.yml/badge.svg?branch=main
[Go Report Card src]: https://goreportcard.com/report/github.com/nafigator/http/reponse/wrapper
[Go Report Card badge]: https://goreportcard.com/badge/github.com/nafigator/http/response/wrapper
[Codecov report src]: https://app.codecov.io/gh/nafigator/http/tree/main
[Codecov report badge]: https://codecov.io/gh/nafigator/http/branch/main/graph/badge.svg
