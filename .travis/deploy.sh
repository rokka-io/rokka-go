#!/bin/sh

set -e

gox -os="linux darwin windows" -arch="amd64" -ldflags "-X github.com/rokka-io/rokka-go/cmd/rokka/cli.cliVersion=${TRAVIS_TAG}" ./cmd/rokka

chmod u+x ./rokka_darwin_amd64 ./rokka_linux_amd64

mkdir build/ && cd build/
cp ../rokka_darwin_amd64 . && mv rokka_darwin_amd64 rokka && zip rokka_macos.zip rokka
cp ../rokka_linux_amd64 . && mv rokka_linux_amd64 rokka && zip rokka_linux.zip rokka
cp ../rokka_windows_amd64.exe . && mv rokka_windows_amd64.exe rokka.exe && zip rokka_windows.zip rokka.exe

cd ..

mv build/*.zip .
rmdir build/