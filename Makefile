GO ?= go
GOFMT ?= gofmt "-s"
GOFILES := $(shell find . -name "*.go")
GOMODULES := $(shell go list ./...)

all:
	$(GO) run cmd/cli/main.go

build:
	$(GO) build -o build/program/app cmd/httpd/main.go

clean:
	@rm -rf build

fmt:
	$(GOFMT) -w $(GOFILES)

test:
	$(GO) clean -testcache
	$(GO) mod tidy
	$(GO) test -cover $(GOMODULES)

update:
	$(GO) get -u ./...

.phony:
	all build clean fmt test update