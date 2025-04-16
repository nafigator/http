<a id="readme-top"></a>
# Go HTTP dumper

[![GitHub release][Release img]][Release src] [![Github main status][Github main status badge]][Github main status src] [![Go Report Card][Go Report Card badge]][Go Report Card src] [![Coverage report][Codecov report badge]][Codecov report src]

Easy to use serverside HTTP dumper.

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#features">Features</a></li>
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
            <li><a href="#error-handling">Error handling</a></li>
            <li><a href="#custom-template">Custom template</a></li>
            <li><a href="#custom-flusher">Custom flusher</a></li>
            <li><a href="#custom-filter">Custom filter</a></li>
            <li><a href="#custom-masker">Custom masker</a></li>
        </ul>
    </li>
    <li><a href="#tests">Tests</a></li>
  </ol>
</details>

## Features
* Easy of use
* Ignore files in body
* Sensitive data masking
* Customizable

## Usage

<details>
  <summary>Example</summary>

```go
package main

import (
  "net/http"

  "github.com/nafigator/http/server/dumper"
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
  d := dumper.New(debug.New(log))
  mux := http.NewServeMux()

  mux.Handle("/", Home)

  srv := http.Server{
    Handler: d.MiddleWare(mux),
  }

  if err := httpServer.ListenAndServe(); err != nil {
    log.Fatal(err)
  }
}

func Home(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Homepage")
}
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>

</details>

After `go run main.go` you'll get output with full HTTP request/response:

<details>
  <summary>Output</summary>

```
2025-01-08 09:18:29.254	DEBUG	HTTP dump:
GET / HTTP/1.1
Host: localhost
Accept-Encoding: gzip



HTTP/1.1 OK 200
Date: Wed, 08 Jan 2025 06:18:29 GMT

Homepage
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>

</details>

> By default, HTTP-requests without *Content-Type* header or *Content-Type* equal `application/octet-stream` dumps
> without HTTP-body.


## Advanced usage
### Masking
Optionally you can mask sensitive data in HTTP dumps using masker. There is 3 masker types:
1. **auth** - masks Authorization header data
2. **query** - masks URL query-params
3. **scalar** - masks JSON scalars in HTTP body.

#### auth

<details>
  <summary>Example</summary>

```go
import (
  "github.com/nafigator/http/server/dumper"
  "github.com/nafigator/http/storage/debug"
  "github.com/nafigator/http/masker/auth"
)

func main() {
  ...
  d := dumper.New(debug.New(log)).
    WithMasker(auth.New()) // Add auth masker
  ...
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>
</details>

This masker will mask dump as follows:

<details>
  <summary>Output</summary>

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
</details>

#### query

<details>
  <summary>Example</summary>

```go
import (
  "github.com/nafigator/http/server/dumper"
  "github.com/nafigator/http/storage/debug"
  "github.com/nafigator/http/masker/query"
)

func main() {
  ...
  d := dumper.New(debug.New(log)).
    WithMasker(query.New([]string{"user","secret"})) // Add query masker
  ...
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>
</details>

This masker will mask dump as follows:

<details>
  <summary>Output</summary>

```
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
</details>

#### json

<details>
  <summary>Example</summary>

```go
import (
  "github.com/nafigator/http/client/dumper"
  "github.com/nafigator/http/storage/debug"
  "github.com/nafigator/http/masker/json"
)

func main() {
  ...
  d := dumper.New(debug.New(log)).
    WithMasker(json.New([]string{"user","secret"})) // Add JSON masker
  ...
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>
</details>

This masker will mask dump as follows:

<details>
  <summary>Output</summary>

```
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
</details>

### Combined masking
Optionally you can combine maskers as follows:

<details>
  <summary>Example</summary>

```go
  ...
  m := auth.New().
    WithNext(json.New([]string{"secret"}))
  d := dumper.New(debug.New(log)).
    WithMasker(m) // Add auth and JSON masker
  ...
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>
</details>

### Control unmasked symbols
By default, all maskers leave 7 unmasked symbols at end for debug purpose. You can change this using `WithUnmasked()`
method.

<details>
  <summary>Example</summary>

```go
  ...
  m := auth.New().WithUnmasked(0)
  d := dumper.New(debug.New(log)).
    WithMasker(m) // Add auth with entire value masker
  ...
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>
</details>

### Error handling
By default, dumper ignores errors from `httputil.DumpRequestOut()` and `httputil.DumpResponse()`. You can change this
by providing error logger with interface:
```go
type logger interface {
  Error(args ...interface{})
}
```

<details>
  <summary>Example</summary>

```go
  ...
  log := zapper.Must(conf.Must())

  d := dumper.
    New(debug.New(log)).
    WithErrLogger(log)
  ...
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>
</details>

### Custom template
You can redefine default dump output template `"HTTP dump:\n%s\n\n%s\n"`. First placeholder is for request, second
for response.

<details>
  <summary>Example</summary>

```go
  ...
  d := dumper.New(debug.New(log)).
    WithTemplate("Dump:\n%s\n✭ ✭ ✭ ✭ ✭ ✭ ✭ ✭ ✭ ✭\n%s\n")
  ...
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>
</details>

This example produces:

<details>
  <summary>Output</summary>

```
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
</details>

### Custom flusher
Flusher contains functionality that provides required message processing. It can be functionality to save
messages to database or output into stdout or anything else. Package [storage/debug][debug src] is an example of 
flusher that uses zap logger under hood to debug messages into stdout. You can implement your own flusher with
interface:
```go
type flusher interface {
  Flush(ctx context.Context, msg string)
}
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Custom filter
If you don't want dump HTTP-bodies requests with specific Content-Type headers, redefine default filter function.
```go
// Requires Content-Type value as param.
func(ct string) bool {
  if ct == mime.Bin || ct == "" {
    return false // do not dump files
  }

  return true
}
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>

<details>
  <summary>Example</summary>

```go
  ...
  d := dumper.New(debug.New(log)).
    WithFilter(func(ct string) bool {
      return ct != mime.PDF // do not dump PDF file in body
    })
  ...
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>
</details>

### Custom masker
You can implement your own masker with interface:
```go
type masker interface {
  Mask(*http.Request, *string)
}
```
Where the second parameter is a pointer to final dump.
<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Tests
Clone repo and run:
```shell
go test
```
<p align="right">(<a href="#readme-top">back to top</a>)</p>

[Release img]: https://img.shields.io/github/v/tag/nafigator/http?logo=github&labelColor=333&color=teal&filter=server/dumper*
[Release src]: https://github.com/nafigator/http/tree/main/server/dumper
[Github main status src]: https://github.com/nafigator/http/tree/main/server/dumper
[Github main status badge]: https://github.com/nafigator/http/actions/workflows/go.yml/badge.svg?branch=main
[Go Report Card src]: https://goreportcard.com/report/github.com/nafigator/http/server/dumper
[Go Report Card badge]: https://goreportcard.com/badge/github.com/nafigator/http/server/dumper
[Codecov report src]: https://app.codecov.io/gh/nafigator/http/tree/main
[Codecov report badge]: https://codecov.io/gh/nafigator/http/branch/main/graph/badge.svg
[debug src]: https://github.com/nafigator/http/tree/main/storage/debug
