BUILD_OUTPUT_DIR=build
GOPATH=$(shell go env GOPATH)
MAIN_BINARY=$(BUILD_OUTPUT_DIR)/tc
VERSION_PKG=main
REPOSITORY=$(shell git config --get remote.origin.url)
REVISION=$(shell git rev-parse HEAD)
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(dir $(mkfile_path))
ORG_PATH="ispringsolutions"
REPO_PATH="${ORG_PATH}/tc"

all: build test doc

.PHONY: ensure
ensure: install_build_tools
	$(GOPATH)/bin/dep ensure -v

.PHONY: build
build: install_build_tools
	$(GOPATH)/bin/dep ensure
	mkdir -p $(BUILD_OUTPUT_DIR)
	CGO_ENABLED=1 go build -v \
	        -ldflags "-X $(VERSION_PKG).REVISION=$(REVISION) -X $(VERSION_PKG).REPOSITORY=$(REPOSITORY)" \
            -o $(MAIN_BINARY) \
            cmd/main.go

.PHONY: clean
clean:
	rm -rf $(BUILD_OUTPUT_DIR)

.PHONY: test
test:
	go test ./... -count=1
	go vet ./...

.PHONY: check
check: install_linter
	$(GOPATH)/bin/golangci-lint run ./... --config .golangci.yml

.PHONY: doc
doc: install_snowboard
	$(GOPATH)/bin/snowboard html -o $(BUILD_OUTPUT_DIR)/tc.html tc.apib

.PHONY: install_build_tools
install_build_tools:
	scripts/install_build_tools.sh

.PHONY: install_linter
install_linter:
	scripts/install_linter.sh

.PHONY: install_snowboard
install_snowboard:
	scripts/install_snowboard.sh

.PHONY: coverage
coverage:
	go test -cover ./...
	go test -coverprofile=build/coverage.out ./...
	go tool cover -func=build/coverage.out
	go tool cover -html=build/coverage.out