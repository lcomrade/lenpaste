NAME = lenpaste
MAIN_GO = ./cmd/*.go

export GOMODULE111=on
LDFLAGS = -w -s -X "main.Version=$(shell git describe --tags --always | sed 's/-/+/' | sed 's/^v//')"


.PHONY: all fmt clean

all:
	mkdir -p ./dist/bin/
	mkdir -p ./dist/share/$(NAME)

	go build -trimpath -ldflags="$(LDFLAGS)" -o ./dist/bin/$(NAME) $(MAIN_GO)
	chmod +x ./dist/bin/$(NAME)

	cp -r ./web ./dist/share/$(NAME)

fmt:
	gofmt -w ./cmd/*.go
	gofmt -w ./internal/apiv1/*.go
	gofmt -w ./internal/config/*.go
	gofmt -w ./internal/logger/*.go
	gofmt -w ./internal/netshare/*.go
	gofmt -w ./internal/raw/*.go
	gofmt -w ./internal/storage/*.go
	gofmt -w ./internal/web/*.go

clean:
	rm -rf ./dist/
