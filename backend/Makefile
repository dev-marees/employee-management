.DEFAULT_GOAL := help
BINARY := bin/api
PKG := ./cmd/api

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

.PHONY: tidy
tidy: ## Sync go.mod / go.sum
	go mod tidy

.PHONY: run
run: ## Run the API locally
	go run $(PKG)

.PHONY: build
build: ## Build the binary into ./bin
	go build -ldflags="-s -w" -o $(BINARY) $(PKG)

.PHONY: test
test: ## Run unit tests
	go test ./... -count=1

.PHONY: test-cover
test-cover: ## Run tests with coverage report
	go test ./... -count=1 -coverprofile=coverage.out
	go tool cover -func=coverage.out

.PHONY: lint
lint: ## Run go vet
	go vet ./...

.PHONY: fmt
fmt: ## Format the code
	gofmt -s -w .

.PHONY: swagger
swagger: ## Generate Swagger docs (requires github.com/swaggo/swag/cmd/swag)
	swag init -g cmd/api/main.go -o docs

.PHONY: docker-build
docker-build: ## Build the Docker image
	docker build -t employee-management-api:latest .

.PHONY: up
up: ## Start the full stack (api + postgres)
	docker compose up --build

.PHONY: down
down: ## Stop the stack and remove volumes
	docker compose down -v
