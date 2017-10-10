# Rokka API client library for Go

*This is currently work in progress*

## Install

```sh
# Install CLI
$ go get github.com/rokka-io/rokka-go/cmd/rokka

# Install library
$ go get github.com/rokka-io/rokka-go/rokka
```

## Development

```sh
# Build during development
$ go install ./cmd/rokka

# Build platform specific executables
$ GOOS=darwin go build -o ./bin/rokka ./cmd/rokka
$ GOOS=linux go build -o ./bin/rokka ./cmd/rokka
$ GOOS=windows go build -o ./bin/rokka ./cmd/rokka
```