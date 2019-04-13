#!/usr/bin/env sh
set -e

export GOPATH=$(go env GOPATH)
export PATH=$PATH:$GOPATH/bin

if ! command -v snowboard>/dev/null; then
    go get -u -v -d github.com/bukalapak/snowboard && \
    (cd "$GOPATH/src/github.com/bukalapak/snowboard" && make install -s)
fi
