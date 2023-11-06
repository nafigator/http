# headers

[![GitHub release][Release img]][Release src] [![Github main status][Github main status badge]][Github main status src] [![Go Report Card][Go Report Card src]][Go Report Card badge] [![Coverage report][Codecov.io report src]][Codecov.io report badge]

HTTP header constants for Gophers.

## Installation

```sh
go install github.com/nafigator/http/headers
```

## Usage

```go
import (
  "fmt"

  "github.com/nafigator/http/headers"
)

fmt.Println(headers.AcceptCharset)
// -> "Accept-Charset"

fmt.Println(headers.IfNoneMatch)
// -> "If-None-Match"

fmt.Println(headers.Normalize("conTent-tYpe"))
// -> "Content-Type"
```

## Tests
```shell
go test -C headers ./...
```


[Release img]: https://img.shields.io/badge/release-1.0.2-green.svg
[Release src]: https://github.com/nafigator/http/headers
[Github main status src]: https://github.com/nafigator/http/tree/main/headers
[Github main status badge]: https://github.com/nafigator/http/actions/workflows/go.yml/badge.svg?branch=main
[Go Report Card src]: https://goreportcard.com/badge/github.com/nafigator/http/headers
[Go Report Card badge]: https://goreportcard.com/report/github.com/nafigator/http/headers
[Codecov.io report src]: https://app.codecov.io/gh/nafigator/http/tree/main
[Codecov.io report badge]: https://codecov.io/gh/nafigator/http/branch/main/graph/badge.svg
