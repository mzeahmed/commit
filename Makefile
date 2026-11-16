# ==============================================================================
# Commit - Development Makefile
# ==============================================================================

.DEFAULT_GOAL := help

GREEN  := \033[0;32m
YELLOW := \033[1;33m
BLUE   := \033[0;34m
RED    := \033[0;31m
RESET  := \033[0m

.PHONY: help build install fmt vet test check tidy \
        release-check release-snapshot clean

help: ## Show available commands
	@echo ""
	@echo "$(BLUE)Commit Development Commands$(RESET)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_-]+:.*##/ {printf "  \033[32m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""

# ==============================================================================
# Development
# ==============================================================================

build: ## Build local binary
	@mkdir -p bin
	go build -o bin/commit .
	@echo "$(GREEN)✓ Binary generated in bin/commit$(RESET)"

install: build ## Build locally and install to /usr/local/bin
	sudo install -m 755 bin/commit /usr/local/bin/commit
	@echo "$(GREEN)✓ Installed in /usr/local/bin/commit$(RESET)"
	@command -v commit

# ==============================================================================
# Quality
# ==============================================================================

fmt: ## Format source code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

test: ## Run unit tests
	go test ./...

check: fmt vet test ## Run all quality checks

# ==============================================================================
# Dependencies
# ==============================================================================

tidy: ## Clean go.mod / go.sum
	go mod tidy

# ==============================================================================
# Release
# ==============================================================================

release-check: ## Validate GoReleaser configuration
	goreleaser check

release-snapshot: ## Simulate a release locally (builds dist/, no GitHub release)
	goreleaser release --snapshot --clean

# ==============================================================================
# Utilities
# ==============================================================================

clean: ## Remove generated files
	rm -rf bin
	rm -rf dist
