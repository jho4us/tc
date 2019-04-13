#!/usr/bin/env sh
 
set -e
 
export GOPATH=$(go env GOPATH)
export PATH=$PATH:$GOPATH/bin
 
if ! command -v golangci-lint>/dev/null; then
    go get -u -v github.com/golangci/golangci-lint/cmd/golangci-lint && \
    cd $(go env GOPATH)/src/github.com/golangci/golangci-lint/cmd/golangci-lint && \
    go install -v -ldflags "-X 'main.version=$(git describe --tags)' -X 'main.commit=$(git rev-parse --short HEAD)' -X 'main.date=$(date)'"
fi