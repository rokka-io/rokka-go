# Rokka API client library for Go

[![GoDoc](https://godoc.org/github.com/rokka-io/rokka-go?status.svg)](https://godoc.org/github.com/rokka-io/rokka-go)
[![Build Status](https://travis-ci.org/rokka-io/rokka-go.svg?branch=master)](https://travis-ci.org/rokka-io/rokka-go)
[![Coverage Status](https://coveralls.io/repos/github/rokka-io/rokka-go/badge.svg?branch=master)](https://coveralls.io/github/rokka-io/rokka-go?branch=master)

*This is currently work in progress*

## Install

Go 1.7 is required.

```sh
# Install CLI
$ go get github.com/rokka-io/rokka-go/cmd/rokka

# Install library
$ go get github.com/rokka-io/rokka-go/rokka
```

## Dependencies

This project uses [dep](https://github.com/golang/dep). Run `dep ensure` for dependencies.

## Development

```sh
# Build during development
$ go install ./cmd/rokka

# Build platform specific executables
$ GOOS=darwin go build -o ./bin/rokka ./cmd/rokka
$ GOOS=linux go build -o ./bin/rokka ./cmd/rokka
$ GOOS=windows go build -o ./bin/rokka ./cmd/rokka

# Update (auto-generate) rokka/operations_objects.go
$ go generate ./rokka
```