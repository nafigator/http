# http
[![GitHub license][License img]][License src] [![Github main status][Github main status badge]][Github main status src] [![Go Report Card][Go Report Card badge]][Go Report Card src] [![Coverage report][Codecov report badge]][Codecov report src] [![Conventional Commits][Conventional commits badge]][Conventional commits src]

Collection of Go HTTP packages

## Packages
### client/dumper
[Package](https://github.com/nafigator/http/blob/main/client/dumper/README.md) for dumping HTTP-client requests/responses.

### masker/auth
[Package](https://github.com/nafigator/http/tree/main/masker/auth) for hiding sensitive data in Authorization header of HTTP-dumps.

### masker/json
[Package](https://github.com/nafigator/http/tree/main/masker/json) for hiding sensitive data in JSON values of HTTP bodies.

### masker/query
[Package](https://github.com/nafigator/http/tree/main/masker/query) for hiding sensitive data in URL-params of HTTP-dumps.

### storage/debug
[Package](https://github.com/nafigator/http/tree/main/storage/debug) with implementation of output HTTP-dumps as logger DEBUG messages.

### headers
[Package](https://github.com/nafigator/http/blob/main/headers/README.md) with constants for HTTP headers.

## Versioning
Each Go module follows *"Semantic Versioning"* specifications. The signature of exported package functions is used
as a public API. Read more on [SemVer.org][semver src].

[License img]: https://img.shields.io/github/license/nafigator/http?color=teal
[License src]: https://www.tldrlegal.com/license/mit-license
[Github main status src]: https://github.com/nafigator/http/tree/main/client/dumper
[Github main status badge]: https://github.com/nafigator/http/actions/workflows/go.yml/badge.svg?branch=main
[Go Report Card src]: https://goreportcard.com/report/github.com/nafigator/http/client/dumper
[Go Report Card badge]: https://goreportcard.com/badge/github.com/nafigator/http/client/dumper
[Codecov report src]: https://app.codecov.io/gh/nafigator/http/tree/main
[Codecov report badge]: https://codecov.io/gh/nafigator/http/branch/main/graph/badge.svg
[Conventional commits src]: https://conventionalcommits.org
[Conventional commits badge]: https://img.shields.io/badge/Conventional%20Commits-1.0.0-teal.svg
[semver src]: http://semver.org
