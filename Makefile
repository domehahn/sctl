.PHONY: build test test-integration lint clean install release-snapshot

BINARY   := sctl
DIST_DIR := dist
MODULE   := github.com/domehahn/sctl
VERSION  ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS  := -ldflags "-X $(MODULE)/internal/cli.Version=$(VERSION) -s -w"

build:
	go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY) ./cmd/sctl

test:
	go test -race -timeout 60s ./...

test-integration:
	go test -race -timeout 120s -tags integration ./tests/integration/...

lint:
	golangci-lint run ./...

clean:
	rm -rf $(DIST_DIR)

install: build
	cp $(DIST_DIR)/$(BINARY) $(GOPATH)/bin/$(BINARY)

release-snapshot:
	goreleaser release --snapshot --clean
