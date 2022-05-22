NAME = lenpaste
MAIN_GO = ./cmd/*.go
LDFLAGS = -w -s

.PHONY: all fmt clean

all:
	mkdir -p ./dist/

	go build -ldflags="$(LDFLAGS)" -o ./dist/$(NAME) $(MAIN_GO)
	chmod +x ./dist/$(NAME)

	cp -r ./web ./dist

fmt:
	gofmt -w ./cmd/*.go
	gofmt -w ./internal/apiv1/*.go
	gofmt -w ./internal/logger/*.go
	gofmt -w ./internal/storage/*.go
	gofmt -w ./internal/web/*.go

clean:
	rm -rf ./dist/
