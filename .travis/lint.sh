#!/bin/sh

set -e

# linting

go get -u honnef.co/go/tools/cmd/megacheck
go get -u golang.org/x/lint/golint
go get -u golang.org/x/tools/cmd/goimports

dirs=$(go list -f {{.Dir}} ./...)
pkgs=$(go list ./...)


echo "Running goimports..."
goimports -d $dirs
echo "Running go tool vet..."
go tool vet $dirs
echo "Running golint ..."
golint $dirs
echo "Running megacheck ..."
megacheck $pkgs