.PHONY: api api-list test vet tools lint analyze check example

GOBIN_DIR := $(or $(shell go env GOBIN),$(shell go env GOPATH)/bin)
GOLANGCI_LINT := $(GOBIN_DIR)/golangci-lint

api:
	@./scripts/list-libgd-api.sh >/dev/null

api-list:
	./scripts/list-libgd-api.sh

test:
	go test -race -count=1 ./...

vet:
	go vet ./...

tools:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

lint: tools
	$(GOLANGCI_LINT) run

analyze: vet lint

check: api vet test

example:
	go run ./examples/basic
