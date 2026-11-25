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
        release-check release-snapshot clean clean-branches

# Branches que clean-branches ne supprime jamais (noms entiers). La branche courante
# est toujours conservée ; HEAD est la référence symbolique du remote.
PROTECTED_BRANCHES := ^(main|develop|HEAD)$$

# Seules les branches fusionnées dans origin/develop sont supprimables : fusion classique
# (la branche est un ancêtre) ou squash (un commit de develop apporte le même correctif).
FILTER_MERGED = while read -r ref name; do \
	if git merge-base --is-ancestor "$$ref" origin/develop \
	|| git cherry origin/develop "$$(git commit-tree "$$ref^{tree}" -p "$$(git merge-base origin/develop "$$ref")" -m squash)" | grep -q '^-'; \
	then echo "$$name"; fi; done

LOCAL_BRANCHES  = git for-each-ref --format='%(refname) %(refname:short)' refs/heads | $(FILTER_MERGED) | grep -vE '$(PROTECTED_BRANCHES)' | grep -vxF "$$(git branch --show-current)"
REMOTE_BRANCHES = git for-each-ref --format='%(refname) %(refname:lstrip=3)' refs/remotes/origin | $(FILTER_MERGED) | grep -vE '$(PROTECTED_BRANCHES)'

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
	sudo install -m 755 bin/commit /usr/local/bin/commit && rm -rf bin/
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

clean-branches: ## Supprimer les branches locales et distantes fusionnées dans develop, sauf main et develop (distantes : selon les règles GitHub)
	@git fetch --prune -q
	@echo "$(YELLOW)Branches locales à supprimer :$(RESET)"
	@out="$$($(LOCAL_BRANCHES))"; [ -n "$$out" ] && echo "$$out" | sed 's/^/  /' || echo "  (aucune)"
	@echo "$(YELLOW)Branches distantes à supprimer :$(RESET)"
	@out="$$($(REMOTE_BRANCHES))"; [ -n "$$out" ] && echo "$$out" | sed 's/^/  /' || echo "  (aucune)"
	@echo ""
	@printf "$(RED)⚠️  Confirmer la suppression ? [y/N] $(RESET)" && read ans && [ "$${ans}" = "y" ] || { echo "$(YELLOW)Annulé.$(RESET)"; exit 1; }

	@echo "$(YELLOW)Suppression des branches locales...$(RESET)"
	@$(LOCAL_BRANCHES) | xargs -r git branch -D || true

	@echo "$(YELLOW)Suppression des branches distantes...$(RESET)"
	@$(REMOTE_BRANCHES) | xargs -r -I {} git push origin --delete {} || true

	@echo "$(GREEN)Nettoyage des branches terminé$(RESET)"
