#!/bin/bash
set -e

if ! [ -d ./release ]; then
  mkdir ./release
fi

VERSION=$(git describe)
GITSHA=$(git rev-parse HEAD)

LDFLAGS="-X main.Version=$VERSION -X main.GitSHA=$GITSHA"
echo $LDFLAGS

go test -v -race .

GOARCH=amd64 go build -ldflags "$LDFLAGS" -o ./release/mesos-cli-v0.0.1-linux-amd64
GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o ./release/mesos-cli-v0.0.1-darwin-amd64
