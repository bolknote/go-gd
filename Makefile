.PHONY: api api-list test vet tools staticcheck lint gosec errcheck ineffassign revive analyze check example

GOBIN_DIR := $(or $(shell go env GOBIN),$(shell go env GOPATH)/bin)
STATICCHECK := $(GOBIN_DIR)/staticcheck
GOLANGCI_LINT := $(GOBIN_DIR)/golangci-lint
GOSEC := $(GOBIN_DIR)/gosec
ERRCHECK := $(GOBIN_DIR)/errcheck
INEFFASSIGN := $(GOBIN_DIR)/ineffassign
REVIVE := $(GOBIN_DIR)/revive

api:
	@./scripts/list-libgd-api.sh >/dev/null

api-list:
	./scripts/list-libgd-api.sh

test:
	go test -race -count=1 ./...

vet:
	go vet ./...

tools:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install github.com/kisielk/errcheck@latest
	go install github.com/gordonklaus/ineffassign@latest
	go install github.com/mgechev/revive@latest

staticcheck: tools
	$(STATICCHECK) ./...

lint: tools
	$(GOLANGCI_LINT) run

gosec: tools
	$(GOSEC) -exclude=G115 ./...

errcheck: tools
	$(ERRCHECK) ./...

ineffassign: tools
	$(INEFFASSIGN) ./...

revive: tools
	$(REVIVE) -config .revive.toml ./...

analyze: vet staticcheck lint gosec errcheck ineffassign revive

check: api vet test

example:
	go run ./examples/basic
