.PHONY: build test lint fmt fmt-check clean help

# Tool definitions
GOLANGCI_LINT = go tool -modfile=tools/go.mod golangci-lint

# Build
build: ## Build the binary
	go build -o bin/crx ./cmd/crx

# Test
test: ## Run tests
	go test -race -shuffle=on ./...

# Lint
lint: ## Run golangci-lint
	$(GOLANGCI_LINT) run --timeout 5m

lint-fix: ## Run golangci-lint with auto-fix
	$(GOLANGCI_LINT) run --fix --timeout 5m

# Format
fmt: ## Format code
	go fmt ./...

fmt-check: ## Check code formatting
	@test -z "$$(gofmt -l .)" || (echo "Files not formatted:"; gofmt -l .; exit 1)

# Clean
clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf dist/

# Help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
