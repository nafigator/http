<a id="readme-top"></a>
# HTTP-client dumper

[![GitHub release][Release img]][Release src] [![Github main status][Github main status badge]][Github main status src] [![Go Report Card][Go Report Card badge]][Go Report Card src] [![Coverage report][Codecov report badge]][Codecov report src]

Easy to use HTTP-client dumper.

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#installation">Installation</a></li>
    <li><a href="#usage">Usage</a></li>
    <li>
        <a href="#advanced-usage">Advanced usage</a>
        <ul>
            <li><a href="#masking">Masking</a></li>
            <ul>
                <li><a href="#auth">auth</a></li>
                <li><a href="#query">query</a></li>
                <li><a href="#json">json</a></li>
            </ul>
            <li><a href="#combined-masking">Combined Masking</a></li>
            <li><a href="#control-unmasked-symbols">Control unmasked symbols</a></li>
            <li><a href="#custom-template">Custom template</a></li>
            <li><a href="#custom-masker">Custom masker</a></li>
        </ul>
    </li>
    <li><a href="#tests">Tests</a></li>
  </ol>
</details>

## Installation

```sh
go install github.com/nafigator/http/client/dumper
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>

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
  d := dumper.New(
    http.DefaultTransport,
    debug.New(log), // Use debug output or implement your own
  )

  c := http.Client{Transport: d}
  _, err := c.Get("https://example.io/api/v3/checks/")
  if err != nil {
    log.Errorln(err)
  }
}
```
After `go run main.go` you'll get output with full HTTP request/response:
```shell
2025-01-08 09:18:29.254	DEBUG	HTTP dump:
GET /api/v3/checks/ HTTP/1.1
Host: example.io
User-Agent: Go-http-client/1.1
Accept-Encoding: gzip



HTTP/2.0 401 Unauthorized
Content-Length: 28
Content-Type: application/json
Date: Wed, 08 Jan 2025 06:18:29 GMT
X-Frame-Options: DENY

{"error": "missing api key"}
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Advanced usage
### Masking
Optionally you can mask sensitive data in HTTP dumps using masker. There is 3 masker types:
1. **auth** - masks Authorization header data
2. **query** - masks URL query-params
3. **scalar** - masks JSON scalars in HTTP body.

#### auth
Example:
```go
import (
  "github.com/nafigator/http/client/dumper"
  "github.com/nafigator/http/storage/debug"
  "github.com/nafigator/http/masker/auth"
)

func main() {
  ...
  // Wrap default http transport by dumper
  d := dumper.New(http.DefaultTransport, debug.New(log)).
	  WithMasker(auth.New()) // Add auth masker
  ...
```

This masker will mask dump as follows:
```shell
2025-01-08 09:18:29.254	DEBUG	HTTP dump:
GET /api/v3/checks/ HTTP/1.1
Host: example.io
Authorization: Bearer ************************f437de0
User-Agent: Go-http-client/1.1
Accept-Encoding: gzip



HTTP/2.0 403 Forbidden
Content-Length: 28
Content-Type: application/json
Date: Wed, 08 Jan 2025 06:18:29 GMT
X-Frame-Options: DENY

{"error": "invalid api key"}
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>

#### query
Example:
```go
import (
  "github.com/nafigator/http/client/dumper"
  "github.com/nafigator/http/storage/debug"
  "github.com/nafigator/http/masker/query"
)

func main() {
  ...
  // Wrap default http transport by dumper
  d := dumper.New(http.DefaultTransport, debug.New(log)).
	  WithMasker(query.New([]string{"user","secret"})) // Add query masker
  ...
```

This masker will mask dump as follows:
```shell
2025-01-08 09:18:29.254	DEBUG	HTTP dump:
GET /api/v3/checks?user=**onymous&secret=*****6789ABC HTTP/1.1
Host: example.io
User-Agent: Go-http-client/1.1
Accept-Encoding: gzip



HTTP/2.0 403 Forbidden
Content-Length: 28
Content-Type: application/json
Date: Wed, 08 Jan 2025 06:18:29 GMT
X-Frame-Options: DENY

{"error": "invalid secret"}
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>

#### json
Example:
```go
import (
  "github.com/nafigator/http/client/dumper"
  "github.com/nafigator/http/storage/debug"
  "github.com/nafigator/http/masker/json"
)

func main() {
  ...
  // Wrap default http transport by dumper
  d := dumper.New(http.DefaultTransport, debug.New(log)).
	  WithMasker(json.New([]string{"user","secret"})) // Add JSON masker
  ...
```

This masker will mask dump as follows:
```shell
2025-01-08 09:18:29.254	DEBUG	HTTP dump:
POST /api/v3/checks/ HTTP/1.1
Host: example.io
User-Agent: Go-http-client/1.1
Content-Type: application/json
Accept-Encoding: gzip

{"user":"**onymous","secret":"*****6789ABC"}


HTTP/2.0 403 Forbidden
Content-Length: 28
Content-Type: application/json
Date: Wed, 08 Jan 2025 06:18:29 GMT
X-Frame-Options: DENY

{"error":"invalid secret"}
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Combined masking
Optionally you can combine maskers as follows:
```go
  ...
  m := auth.New().
    WithNext(json.New([]string{"secret"}))
  // Wrap default http transport by dumper
  d := dumper.New(http.DefaultTransport, debug.New(log)).
    WithMasker(m) // Add auth and JSON masker
  ...
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Control unmasked symbols
By default, all maskers leave 7 unmasked symbols at end for debug purpose. You can change this using `WithUnmasked()`
method. Example:
```go
  ...
  m := auth.New().WithUnmasked(0)
  // Wrap default http transport by dumper
  d := dumper.New(http.DefaultTransport, debug.New(log)).
    WithMasker(m) // Add auth with entire value masker
  ...
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Custom template
You can redefine default dump output template `"HTTP dump:\n%s\n\n%s\n"`. First placeholder is for request, second
for response. Example:
```go
  ...
  // Wrap default http transport by dumper
  d := dumper.New(
    http.DefaultTransport,
    debug.New(log), // Use debug output or implement your own
  )
  d.WithTemplate("Dump:\n%s\n✭ ✭ ✭ ✭ ✭ ✭ ✭ ✭ ✭ ✭\n%s\n")
  ...
```
This example produces output:
```shell
2025-01-09 16:03:20.461	DEBUG	Dump:
GET /api/v3/checks/ HTTP/1.1
Host: example.io
User-Agent: Go-http-client/1.1
Accept-Encoding: gzip


✭ ✭ ✭ ✭ ✭ ✭ ✭ ✭ ✭ ✭

HTTP/2.0 401 Unauthorized
Content-Length: 28
Content-Type: application/json
Date: Thu, 09 Jan 2025 13:03:20 GMT
X-Frame-Options: DENY

{"error": "missing api key"}
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Custom masker
You can implement your own masker with interface:
```go
  type masker interface {
    Mask(*http.Request, *string)
  }
```
Where second parameter is pointer to final dump.
<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Tests
Clone repo and run:
```shell
go test -C tests ./...
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>

[Release img]: https://img.shields.io/github/v/tag/nafigator/http?logo=github&labelColor=333&color=teal&filter=client/dumper*
[Release src]: https://github.com/nafigator/http/tree/main/client/dumper
[Github main status src]: https://github.com/nafigator/http/tree/main/client/dumper
[Github main status badge]: https://github.com/nafigator/http/actions/workflows/go.yml/badge.svg?branch=main
[Go Report Card src]: https://goreportcard.com/report/github.com/nafigator/http/client/dumper
[Go Report Card badge]: https://goreportcard.com/badge/github.com/nafigator/http/client/dumper
[Codecov report src]: https://app.codecov.io/gh/nafigator/http/tree/main
[Codecov report badge]: https://codecov.io/gh/nafigator/http/branch/main/graph/badge.svg
