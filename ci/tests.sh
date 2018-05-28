#!/bin/bash
set -eu

mkdir -p "$GOPATH/src/github.com/pivotal-cloudops"
cp -R git-omen "$GOPATH/src/github.com/pivotal-cloudops/omen"
go get -u github.com/golang/dep/cmd/dep
go get github.com/onsi/ginkgo/...


OUTPUT_DIR=$(pwd)

cd "$GOPATH/src/github.com/pivotal-cloudops/omen"

dep ensure
ginkgo -r
go build -o ${OUTPUT_DIR}/omen-build/omen-linux64
