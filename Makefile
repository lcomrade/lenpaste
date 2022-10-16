GO ?= go
GOFMT ?= gofmt

NAME = lenpaste
MAIN_GO = ./cmd/*.go

export GOMODULE111=on
LDFLAGS = -w -s -X "main.Version=$(shell git describe --tags --always | sed 's/-/+/' | sed 's/^v//')"


.PHONY: all fmt clean

all:
	mkdir -p ./dist/bin/

	$(GO) build -trimpath -ldflags="$(LDFLAGS)" -o ./dist/bin/$(NAME) $(MAIN_GO)
	chmod +x ./dist/bin/$(NAME)

fmt:
	@$(GO) -w $(shell find ./ -type f -name '*.go')

clean:
	rm -rf ./dist/
