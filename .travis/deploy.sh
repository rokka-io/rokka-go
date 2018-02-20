#!/bin/sh

set -e

gox -os="linux darwin windows" -arch="amd64" -ldflags "-X github.com/rokka-io/rokka-go/cmd/rokka/cli.cliVersion=${TRAVIS_TAG}" ./cmd/rokka

chmod u+x ./rokka_darwin_amd64 ./rokka_linux_amd64