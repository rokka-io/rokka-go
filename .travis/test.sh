#!/bin/sh

set -e

# linting
goimports -d ./
go tool vet ./

# testing
go test -v ./...

# coveralls.io (see https://github.com/golang/go/issues/6909 why the use of overalls)
overalls -project=github.com/rokka-io/rokka-go -covermode=count
goveralls -coverprofile=overalls.coverprofile -service=travis-ci