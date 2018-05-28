#!/bin/bash
set -eu

mkdir -p "$GOPATH/src/github.com/pivotal-cloudops"
cp -R git-omen "$GOPATH/src/github.com/pivotal-cloudops/omen"
go get -u github.com/golang/dep/cmd/dep

cd "$GOPATH/src/github.com/pivotal-cloudops/omen"

dep ensure
go test -v ./...
