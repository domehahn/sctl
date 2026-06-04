.PHONY: build test test-integration test-all coverage lint clean install release-snapshot

BINARY   := skpm
DIST_DIR := dist
MODULE   := github.com/domehahn/sctl
VERSION  ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS  := -ldflags "-X $(MODULE)/internal/cli.Version=$(VERSION) -s -w"
INSTALL_DIR ?= $(shell gobin="$$(go env GOBIN)"; if [ -n "$$gobin" ]; then echo "$$gobin"; else echo "$$(go env GOPATH)/bin"; fi)
UNIT_TEST_PKGS := ./tests/unit/...
INTEGRATION_TEST_PKGS := ./tests/integration/...
ALL_TEST_PKGS := ./tests/...

build:
	go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY) ./cmd/skpm

test:
	go test -race -timeout 60s $(UNIT_TEST_PKGS)

test-integration:
	go test -race -timeout 120s -tags integration $(INTEGRATION_TEST_PKGS)

test-all:
	go test -race -timeout 120s $(ALL_TEST_PKGS)

coverage:
	go test -coverpkg=./internal/... -coverprofile=coverage.out $(ALL_TEST_PKGS)
	go tool cover -func=coverage.out

lint:
	golangci-lint run ./...

clean:
	rm -rf $(DIST_DIR)

install: build
	mkdir -p $(INSTALL_DIR)
	cp $(DIST_DIR)/$(BINARY) $(INSTALL_DIR)/$(BINARY)
	@echo "Installed $(BINARY) to $(INSTALL_DIR)/$(BINARY)"

release-snapshot:
	goreleaser release --snapshot --clean
