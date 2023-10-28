# headers

[![GitHub release][Release img]][Release src] [![Github main status][Github main status badge]][Github main status src]

HTTP header constants for Gophers.

## Installation

```sh
go install github.com/nafigator/http/headers
```

## Documentation

https://godoc.org/github.com/nafigator/http/headers

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

[Release img]: https://img.shields.io/badge/release-0.1.0-red.svg
[Release src]: https://github.com/nafigator/http/headers
[Github main status src]: https://github.com/nafigator/http/tree/main/headers
[Github main status badge]: https://github.com/nafigator/http/actions/workflows/go.yml/badge.svg?branch=feat
