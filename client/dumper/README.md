# HTTP-client dumper

[![GitHub release][Release img]][Release src] [![Github main status][Github main status badge]][Github main status src] [![Go Report Card][Go Report Card badge]][Go Report Card src] [![Coverage report][Codecov report badge]][Codecov report src]

Easy to use HTTP-client dumper.

## Installation

```sh
go install github.com/nafigator/http/client/dumper
```

## Usage

```go
package main

import (
  "net/http"

  "github.com/nafigator/http/client/dumper"
  "github.com/nafigator/http/storage/debug"
  "github.com/nafigator/zapper"
  "github.com/nafigator/zapper/conf"
)

const (
  zapConfig = `
level: debug
encoding: console
outputPaths:
  - stdout
errorOutputPaths:
  - stderr
encoderConfig:
  messageKey: message
  levelKey:   level
  timeKey:    time
  levelEncoder: capital
  timeEncoder:
    layout: 2006-01-02 15:04:05.000
`
)

func main() {
  log := zapper.Must(conf.MustYML(zapConfig))

  // Wrap default http transport by dumper
  t := dumper.New(
    "",
    http.DefaultTransport,
    nil,
    debug.New(log), // Use debug output
    log,
  )

  c := http.Client{Transport: t}
  _, err := c.Get("https://healthchecks.io/api/v3/checks/")
  if err != nil {
    log.Errorln(err)
  }
}
```
After `go run main.go` you'll get output with full HTTP request/response:
```shell
2025-01-08 09:18:29.254	DEBUG	HTTP dump:
GET /api/v3/checks/ HTTP/1.1
Host: healthchecks.io
User-Agent: Go-http-client/1.1
Accept-Encoding: gzip



HTTP/2.0 401 Unauthorized
Content-Length: 28
Access-Control-Allow-Headers: X-Api-Key
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Origin: *
Access-Control-Max-Age: 600
Content-Type: application/json
Cross-Origin-Opener-Policy: same-origin
Date: Wed, 08 Jan 2025 06:18:29 GMT
Referrer-Policy: strict-origin-when-cross-origin
Server: nginx
Vary: Cookie
X-Content-Type-Options: nosniff
X-Frame-Options: DENY

{"error": "missing api key"}
```

## Tests
Clone repo and run:
```shell
go test -C tests ./...
```

[Release img]: https://img.shields.io/github/v/tag/nafigator/http?logo=github&labelColor=333&color=teal&filter=client/dumper*
[Release src]: https://github.com/nafigator/http/tree/main/client/dumper
[Github main status src]: https://github.com/nafigator/http/tree/main/client/dumper
[Github main status badge]: https://github.com/nafigator/http/actions/workflows/go.yml/badge.svg?branch=main
[Go Report Card src]: https://goreportcard.com/report/github.com/nafigator/http/client/dumper
[Go Report Card badge]: https://goreportcard.com/badge/github.com/nafigator/http/client/dumper
[Codecov report src]: https://app.codecov.io/gh/nafigator/http/tree/main
[Codecov report badge]: https://codecov.io/gh/nafigator/http/branch/main/graph/badge.svg
