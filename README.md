<img src="http://vinci-proxy.github.io/public/images/large.png" />

[![Build Status](https://travis-ci.org/vinci-proxy/vinci.png)](https://travis-ci.org/vinci-proxy/vinci) [![GitHub release](https://img.shields.io/badge/version-0.1.0-orange.svg?style=flat)](https://github.com/vinci-proxy/vinci/releases) [![GoDoc](https://godoc.org/github.com/vinci-proxy/vinci?status.svg)](https://godoc.org/github.com/vinci-proxy/vinci) [![Coverage Status](https://coveralls.io/repos/github/vinci-proxy/vinci/badge.svg?branch=master)](https://coveralls.io/github/vinci-proxy/vinci?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/vinci-proxy/vinci)](https://goreportcard.com/report/github.com/vinci-proxy/vinci) [![Status](https://img.shields.io/badge/status-beta-blue.svg)](#)

**Note**: vinci is still beta under heavy development.

## Goals

- Simple, idiomatic, hackable API.
- Designed for strong composability and extensibility.
- Fully middleware-oriented.
- Built on top of `net/http`.
- Built-in (but optional) HTTP router and multiplexer.
- Great interpolarity with standard `net/http` interfaces
- Great convergence with third-party HTTP interfaces (Gorilla, Negroni, Alice...).
- Small core and code base.

## Installation

```bash
go get -u gopkg.in/vinci-proxy/vinci.v0
```

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

## API

See [godoc reference](https://godoc.org/github.com/vinci-proxy/vinci) for detailed API documentation.

## Examples

See [examples](https://github.com/vinci-proxy/vinci/tree/master/_examples) directory.

## License

[Apache License](https://opensource.org/licenses/Apache-2.0) 2.0 and [MIT License](https://opensource.org/licenses/MIT).
