#!/bin/bash
set -eu

mkdir -p "$GOPATH/src/github.com/pivotal-cloudops"
cp -R omen "$GOPATH/src/github.com/pivotal-cloudops/"

cd "$GOPATH/src/github.com/pivotal-cloudops/omen"
go get -u github.com/golang/dep/cmd/dep

dep ensure
go test -v ./...
