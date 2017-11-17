#!/bin/sh

set -e

gox -os="linux darwin windows" -arch="amd64" -ldflags "-X main.cliVersion=${TRAVIS_TAG}" ./cmd/rokka