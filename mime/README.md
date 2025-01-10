# MIME

[![GitHub release][Release img]][Release src] [![Github main status][Github main status badge]][Github main status src] [![Go Report Card][Go Report Card badge]][Go Report Card src] [![Coverage report][Codecov report badge]][Codecov report src]

MIME types constants for Gophers.

## Installation

```sh
go install github.com/nafigator/http/mime
```

## Usage

```go
import (
  "fmt"

  "github.com/nafigator/http/mime"
)

fmt.Println(mime.JSON)   // application/json
fmt.Println(mime.PDF)    // application/pdf
fmt.Println(mime.Bin)    // application/octet-stream
```

## Intention
Minimalistic constants for elimination in code mess like:
```go
  ...
  // w.Header().Set("Content-Type", mime.XLSX)
  w.Header().Set("Content-Type", mimetype.ApplicationVndOpenxmlformatsOfficedocumentSpreadsheetmlSheet)
  ...
```

[Release img]: https://img.shields.io/github/v/tag/nafigator/http?logo=github&labelColor=333&color=teal&filter=mime*
[Release src]: https://github.com/nafigator/http/tree/main/mime
[Github main status src]: https://github.com/nafigator/http/tree/main/mime
[Github main status badge]: https://github.com/nafigator/http/actions/workflows/go.yml/badge.svg?branch=main
[Go Report Card src]: https://goreportcard.com/report/github.com/nafigator/http/mime
[Go Report Card badge]: https://goreportcard.com/badge/github.com/nafigator/http/mime
[Codecov report src]: https://app.codecov.io/gh/nafigator/http/tree/main
[Codecov report badge]: https://codecov.io/gh/nafigator/http/branch/main/graph/badge.svg
