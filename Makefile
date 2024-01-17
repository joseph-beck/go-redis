SWAG ?= swag
GO ?= go
GOFMT ?= gofmt "-s"
GOFILES := $(shell find . -name "*.go")
GOMODULES := $(shell go list ./...)

all:
	$(GO) run main.go

build:
	$(GO) build -o build/program/app main.go

clean:
	@rm -rf build
	$(GO) clean

docs:
	$(SWAG) i --dir ./,\
	./services/

fmt:
	$(GOFMT) -w $(GOFILES)

test:
	$(GO) clean -testcache
	$(GO) mod tidy
	$(GO) test -cover $(GOMODULES)

update:
	$(GO) get -u ./...
	$(GO) mod tidy

info:
	@$(GO) vet $(GOMODULES)
	@$(GO) list $(GOMODULES)
	@$(GO) version

.phony:
	all build clean fmt test update info
