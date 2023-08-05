.PHONY: build
build: ## Build the project binary
	go build -o gomodtrace ./cmd/gomodtrace

.PHONY: test
test: ## Run project tests
	go test ./...

.PHONY: integration_test
integration_test: ## Run integration tests
	# success case
	(echo A B | go run ./cmd/gomodtrace A B) || (echo "E: success case failed" && exit 255)
	# fail when no input provided
	go run ./cmd/gomodtrace; test $$? -eq 1 || (echo "E: test fail on no input provided" && exit 255)
	# fail when no arguments provided
	(echo A B | go run ./cmd/gomodtrace); test $$? -eq 1 || (echo "E: test fail on no arguments provided" && exit 255)
	#
	@echo "Test success"

.PHONY: check
check: test integration_test ## Run all tests


.PHONY: help
help: ## Print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'