#!/bin/sh

set -e

# linting

go get -u honnef.co/go/tools/cmd/megacheck
go get -u golang.org/x/lint/golint
go get -u golang.org/x/tools/cmd/goimports

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