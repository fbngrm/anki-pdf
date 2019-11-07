.PHONY: all build test test-race test-cover

all: build

build:
	mkdir -p bin
	go build -o bin/anki-pdf \
        -ldflags "-X main.version=$${VERSION:-$$(git describe --tags --always --dirty)}" \
        ./cmd/anki-pdf/main.go
