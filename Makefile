BINARY_NAME = "gwid-core"
BINARY_PATH=app/$(BINARY_NAME)
MAIN_PATH = "./cmd/main.go"
TMP_DIR = "tmp"
DOCKER_IMAGE = BINARY_NAME

.DEFAULT_GOAL := help

.PHONY: build
build:
	@echo "Building application..."
	@mkdir -p app
	@go build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Binary built: $(BINARY_PATH)"

.PHONY: build-dev
build-dev:
	@echo "Building application for development..."
	@mkdir -p tmp
	@go build -gcflags="all=-N -l" -o ./tmp/main $(MAIN_PATH)

.PHONY: run
run: build
	@echo "Running application..."
	@./$(BINARY_PATH)

.PHONY: dev
dev: deps
	@echo "Starting development server with hot-reloading..."
	@air -c .air.toml

.PHONY: docker-dev
docker-dev:
	@echo "(docker) Starting development server with hot hot-reloading..."
	@docker compose -f dev-docker-compose.yml up


.PHONY: docker-prod
docker-prod:
	@echo "(docker) Starting production server..." 
	@docker compose -f docker-compose.yml up

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/ tmp/ app/ build-errors.log
	@go clean

.PHONY: deps
deps:
	@echo "Tidying modules..."
	@go mod tidy
	@echo "Tidy complete!"

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Help target
.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
