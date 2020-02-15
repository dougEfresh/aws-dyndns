BIN_DIR := .tools/bin

GO := go
ifdef GO_BIN
	GO = $(GO_BIN)
endif

GOLANGCI_LINT_VERSION := 1.21.0
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint
GIT_COMMIT := $(shell git rev-parse --short HEAD 2> /dev/null || echo "no-revision")
GIT_COMMIT_MESSAGE := $(shell git show -s --format='%s' 2> /dev/null | tr ' ' _ | tr -d "'")
GIT_TAG := $(shell git describe --tags 2> /dev/null || echo "no-tag")
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2> /dev/null || echo "no-branch")
BUILD_TIME := $(shell date +%FT%T%z)
GOFILES := $(shell find . -name '*.go' -type f)

## all: The default target. Build, test, lint
all: test lint

## tidy: go mod tidy
tidy:
	$(GO) mod tidy -v

## fmt: format all go code
fmt:
	gofmt -s -w .

## build: build all files, including protoc if included
build: test
	go build ./...

install: aws-dyndns
	 install -m755 -D aws-dyndns $(DESTDIR)/usr/bin/aws-dyndns

## test: Run all tests
test:
	$(GO) test -cover -race -v ./...

## test-coverate: Run all tests and collect coverage
test-coverage:
	$(GO) test -race ./... -v -covermode=count -coverprofile=testCoverage.txt

## lint: lint all go code
lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run --fast --enable-all -D wsl

$(GOLANGCI_LINT):
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN_DIR) v$(GOLANGCI_LINT_VERSION)

aws-dyndns: $(GOFILES)
		$(GO) build -o aws-dyndns -v -ldflags '-X $(VERSION_PACKAGE).GitHash=$(GIT_COMMIT) -X $(VERSION_PACKAGE).GitTag=$(GIT_TAG) -X $(VERSION_PACKAGE).GitBranch=$(GIT_BRANCH) -X $(VERSION_PACKAGE).BuildTime=$(BUILD_TIME) -X $(VERSION_PACKAGE).GitCommitMessage=$(GIT_COMMIT_MESSAGE)' main.go
