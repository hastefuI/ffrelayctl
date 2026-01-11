.PHONY: setup
setup: ## Set up development environment
	@echo "Checking dependencies..."
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "Error: Docker is not installed or not in PATH"; \
		echo "The git hooks setup requires Docker to run commitlint"; \
		echo "Please install Docker: https://docs.docker.com/get-started/get-docker/"; \
		exit 1; \
	fi
	@echo "Docker found"
	@echo "Configuring git hooks..."
	@git config core.hooksPath .githooks
	@chmod +x .githooks/commit-msg
	@echo "Git hooks installed"

.PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
