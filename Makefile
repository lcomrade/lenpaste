GO ?= go
GOFMT ?= gofmt

NAME = lenpaste
VERSION = $(shell git describe --tags --always | sed 's/-/+/' | sed 's/^v//')

MAIN_GO = ./cmd/$(NAME)/*.go
LDFLAGS = -w -s -X "main.Version=$(VERSION)"

.PHONY: all tarball fmt clean

all:
	mkdir -p ./dist/bin/

	$(GO) build -trimpath -ldflags="$(LDFLAGS)" -o ./dist/bin/$(NAME) $(MAIN_GO)
	chmod +x ./dist/bin/$(NAME)

tarball:
	mkdir -p ./dist/tmp/$(NAME)-$(VERSION)/
	cp -r $(filter-out ./. ./.. ./.git ./dist,$(shell echo ./* ./.*)) ./dist/tmp/$(NAME)-$(VERSION)/

	go mod vendor -o ./dist/tmp/$(NAME)-$(VERSION)/vendor/

	sed -i '/^COPY .*go.mod/d'        ./dist/tmp/$(NAME)-$(VERSION)/Dockerfile
	sed -i '/^RUN go mod download/d'  ./dist/tmp/$(NAME)-$(VERSION)/Dockerfile
	sed -i "s/ git / /"               ./dist/tmp/$(NAME)-$(VERSION)/Dockerfile
	sed -i "s/ ca-certificates / /"   ./dist/tmp/$(NAME)-$(VERSION)/Dockerfile

	sed -i "/^VERSION[[:space:]]*=/c\VERSION=$(VERSION)"  ./dist/tmp/$(NAME)-$(VERSION)/Makefile
	sed -i "s/\$$(GO) build/\$$(GO) build -mod=vendor/"  ./dist/tmp/$(NAME)-$(VERSION)/Makefile

	mkdir -p ./dist/sources/
	tar --mtime="$(git log -1 --format=%ai)" \
		-C ./dist/tmp/ \
		-zcf ./dist/sources/$(NAME)-$(VERSION).tar.gz $(NAME)-$(VERSION)/
	rm -rf ./dist/tmp/

fmt:
	@$(GOFMT) -w $(shell find ./ -type f -name '*.go')

clean:
	rm -rf ./dist/
