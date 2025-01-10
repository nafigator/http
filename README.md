<a id="readme-top"></a>
# http
[![GitHub license][License img]][License src] [![Github main status][Github main status badge]][Github main status src] [![Coverage report][Codecov report badge]][Codecov report src] [![Conventional Commits][Conventional commits badge]][Conventional commits src] ![Repo size][Repo size badge]

Collection of Go HTTP packages

### Packages
#### client/dumper
[Package](https://github.com/nafigator/http/blob/main/client/dumper/README.md) for dumping HTTP-client requests/responses.

#### masker/auth
[Package](https://github.com/nafigator/http/tree/main/masker/auth) for hiding sensitive data in Authorization header of HTTP-dumps.

#### masker/json
[Package](https://github.com/nafigator/http/tree/main/masker/json) for hiding sensitive data in JSON values of HTTP bodies.

#### masker/query
[Package](https://github.com/nafigator/http/tree/main/masker/query) for hiding sensitive data in URL-params of HTTP-dumps.

#### storage/debug
[Package](https://github.com/nafigator/http/tree/main/storage/debug) with implementation of output HTTP-dumps as logger DEBUG messages.

#### headers
[Package](https://github.com/nafigator/http/blob/main/headers/README.md) with constants for HTTP headers.

#### mime
[Package](https://github.com/nafigator/http/blob/main/mime/README.md) with constants for MIME types.

### Versioning
Each Go module follows *"Semantic Versioning"* specifications. The signature of exported package functions is used
as a public API. Read more on [SemVer.org][semver src].

### Built with

[![Go][Go badge]][Go URL]&nbsp;&nbsp;&nbsp;&nbsp;[![GoLand][GoLand badge]][GoLand URL]&nbsp;&nbsp;&nbsp;&nbsp;[![Git][Git badge]][Git URL]

[![Docker][Docker badge]][Docker URL]&nbsp;&nbsp;&nbsp;&nbsp;[![Codecov][Codecov badge]][Codecov URL]&nbsp;&nbsp;&nbsp;&nbsp;[![GitHub Actions][Github actions badge]][Github actions URL]

<p align="right">(<a href="#readme-top">back to top</a>)</p>
</details>

[License img]: https://img.shields.io/github/license/nafigator/http?color=teal
[License src]: https://www.tldrlegal.com/license/mit-license
[Github main status src]: https://github.com/nafigator/http/tree/main
[Github main status badge]: https://github.com/nafigator/http/actions/workflows/go.yml/badge.svg?branch=main
[Codecov report src]: https://app.codecov.io/gh/nafigator/http/tree/main
[Codecov report badge]: https://codecov.io/gh/nafigator/http/branch/main/graph/badge.svg
[Conventional commits src]: https://conventionalcommits.org
[Conventional commits badge]: https://img.shields.io/badge/Conventional%20Commits-1.0.0-teal.svg
[Repo size badge]: https://img.shields.io/github/repo-size/nafigator/http?logo=github&color=teal
[Go badge]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=fff&logoSize=auto
[Go URL]: https://go.dev
[GoLand badge]: https://img.shields.io/badge/GoLand-000?&style=for-the-badge&logo=goland&logoColor=FF318C&logoSize=auto
[GoLand URL]: https://www.jetbrains.com/go/
[Git badge]: https://img.shields.io/badge/Git-fff?style=for-the-badge&logo=git&logoColor=F05032
[Git URL]: https://git-scm.com/
[Docker badge]: https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=fff
[Docker URL]: https://www.docker.com/
[Codecov badge]: https://img.shields.io/badge/codecov-ff0077?style=for-the-badge&logo=codecov&logoColor=fff
[Codecov URL]: https://codecov.io/
[Github actions badge]: https://img.shields.io/badge/GitHub%20Actions-2088FF?style=for-the-badge&logo=githubactions&logoColor=fff&logoSize=auto&labelColor=githubactions
[Github actions URL]: https://github.com/nafigator/http/actions
[semver src]: http://semver.org
