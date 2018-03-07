#!/bin/sh

set -e

# linting

# go 1.8 doesn't exclude /vendor/ yet.
dirs=$(go list -f {{.Dir}} ./... | grep -v /vendor/)
pkgs=$(go list ./... | grep -v /vendor/)

echo "Running goimports..."
goimports -d $dirs
echo "Running go tool vet..."
go tool vet $dirs
echo "Running golint ..."
golint $dirs
echo "Running megacheck ..."
megacheck $pkgs

# testing

echo "Running tests..."
go test $pkgs

# race detection
echo "Running tests with race detector enabled..."
go test -race -count 3 $pkgs

# coveralls.io (see https://github.com/golang/go/issues/6909 why the use of overalls)
echo "Running coveralls..."
overalls -project=github.com/rokka-io/rokka-go -covermode=count
goveralls -coverprofile=overalls.coverprofile -service=travis-ci