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
	@gofmt -w $(shell find ./ -type f -name '*.go')

clean:
	rm -rf ./dist/
