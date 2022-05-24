NAME = lenpaste
MAIN_GO = ./cmd/*.go
LDFLAGS = -w -s

.PHONY: all fmt clean

all:
	mkdir -p ./dist/bin/
	mkdir -p ./dist/share/$(NAME)

	go build -ldflags="$(LDFLAGS)" -o ./dist/bin/$(NAME) $(MAIN_GO)
	chmod +x ./dist/bin/$(NAME)

	cp -r ./web ./dist/share/$(NAME)

fmt:
	gofmt -w ./cmd/*.go
	gofmt -w ./internal/apiv1/*.go
	gofmt -w ./internal/logger/*.go
	gofmt -w ./internal/netshare/*.go
	gofmt -w ./internal/storage/*.go
	gofmt -w ./internal/web/*.go

clean:
	rm -rf ./dist/
