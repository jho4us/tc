#!/usr/bin/env sh

set -e
export GOPATH=$(go env GOPATH)
export PATH=$PATH:$GOPATH/bin

if ! command -v dep>/dev/null; then
    go get -u -v -d github.com/golang/dep/cmd/dep
    go install -v github.com/golang/dep/cmd/dep
fi