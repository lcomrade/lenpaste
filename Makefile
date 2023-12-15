GO ?= go
GOFMT ?= gofmt

NAME = lenpaste
VERSION = $(shell git describe --tags --always | sed 's/-/+/' | sed 's/^v//')

CGO_ENABLED ?= 1
GOOS ?=
GOARCH ?=
MAIN_GO = ./cmd/$(NAME)/*.go
LDFLAGS = -w -s -X "main.Version=$(VERSION)"

.PHONY: all fmt clean

all:
	mkdir -p ./dist/bin/

	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) $(GO) build -trimpath -ldflags="$(LDFLAGS)" -o ./dist/bin/$(NAME) $(MAIN_GO)
	chmod +x ./dist/bin/$(NAME)

fmt:
	@$(GOFMT) -w $(shell find ./ -type f -name '*.go')

clean:
	rm -rf ./dist/
