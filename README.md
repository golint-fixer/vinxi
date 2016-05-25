<img src="http://vinxi.github.io/public/images/large.png" />

[![Build Status](https://travis-ci.org/vinxi/vinxi.png)](https://travis-ci.org/vinxi/vinxi) [![GitHub release](https://img.shields.io/badge/version-0.1.0-orange.svg?style=flat)](https://github.com/vinxi/vinxi/releases) [![GoDoc](https://godoc.org/github.com/vinxi/vinxi?status.svg)](https://godoc.org/github.com/vinxi/vinxi) [![Coverage Status](https://coveralls.io/repos/github/vinxi/vinxi/badge.svg?branch=master)](https://coveralls.io/github/vinxi/vinxi?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/vinxi/vinxi)](https://goreportcard.com/report/github.com/vinxi/vinxi) [![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/vinxi/vinxin/blob/master/LICENSE.md) [![Status](https://img.shields.io/badge/status-beta-blue.svg)](#) 

**Note**: vinxi is still beta under heavy development.

## Features

- Strong modularity via hirarchical middleware layer.
- Multiple middleware interfaces supported (e.g: http.Handler, Negroni, Alice...)
- Built-in multiplexer for easy composition.
- Idiomatic built on top of `net/http` package.
- Default HTTP/S and WebSocket traffic forward.
- Built-in middleware components (e.g: logging, metrics, service discovery, balancer...).
- Tiny, hackable core.
- Completely written in Go. No dependencies.

## Installation

```bash
go get -u gopkg.in/vinxi/vinxi.v0
```

<!--
## Docs

- Introduction
- Installation
- API
- Design goals
- Use cases
- Middleware layer
- List of middleware
- Interpolarity with existent frameworks and libraries.
- Writting a middleware
- Performance
- Benchmarking
- Examples
-->

## API

See [godoc reference](https://godoc.org/github.com/vinxi/vinxi) for detailed API documentation.

## Examples

See [examples](https://github.com/vinxi/vinxi/tree/master/_examples) directory.

## Command-line interface

See [vinxictl](https://github.com/vinxi/vinxictl) for command-line usage.

## Development

Clone the repository:
```bash
git clone https://github.com/vinxi/vinxi.git && cd vinxi
```

Create subpackages symbolic links in `$GOPATH`:
```bash
make link
```

Lint, format and run tests:
```bash
make 
```

## License

Mixed [Apache License](https://opensource.org/licenses/Apache-2.0) 2.0 and [MIT License](https://opensource.org/licenses/MIT) (see file header for details).
