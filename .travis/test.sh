#!/bin/sh

set -e

# linting
echo "Running goimports..."
goimports -d $(go list -f {{.Dir}} ./... | grep -v /vendor/)
echo "Running go tool vet..."
go tool vet $(go list -f {{.Dir}} ./... | grep -v /vendor/)

# testing
echo "Starting tests..."
go test $(go list ./...)

# race detection
echo "Running tests with race detector enabled..."
go test -race -count 3 $(go list ./...)

# coveralls.io (see https://github.com/golang/go/issues/6909 why the use of overalls)
echo "Running coveralls..."
overalls -project=github.com/rokka-io/rokka-go -covermode=count
goveralls -coverprofile=overalls.coverprofile -service=travis-ci