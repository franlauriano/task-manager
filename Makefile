# Variables
ETC_DIR := etc
VAR_DIR := var
BIN_DIR := bin
BIN_NAME := task-api
ENV_FILE := $(ETC_DIR)/.env
TEST_ENV_FILE := $(ETC_DIR)/.env.test
CONFIG_FILE := $(ETC_DIR)/config.toml
CONFIG_EXAMPLE := $(ETC_DIR)/config.toml.example
ENV_EXAMPLE := $(ETC_DIR)/.env.example

# Auto-initialize configuration files
$(shell if [ ! -f $(ENV_FILE) ] && [ -f $(ENV_EXAMPLE) ]; then cp $(ENV_EXAMPLE) $(ENV_FILE); fi)
$(shell if [ ! -f $(CONFIG_FILE) ] && [ -f $(CONFIG_EXAMPLE) ]; then cp $(CONFIG_EXAMPLE) $(CONFIG_FILE); fi)

# Load environment variables
ifneq (,$(wildcard $(ENV_FILE)))
    include $(ENV_FILE)
    export
endif

.PHONY: help deps build clean db-up db-down migrate migrate-down seed run run-dev test coverage

# Helper to validate file
check-env = @if [ ! -f $(1) ]; then echo "Error: $(1) not found"; exit 1; fi

help: ## Show help
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Download and organize dependencies
	go mod download
	go mod tidy

build: deps ## Build binary
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/${BIN_NAME} cmd/main.go

clean: ## Remove generated files
	@rm -rf $(BIN_DIR)/${BIN_NAME} $(VAR_DIR)/coverage.out $(VAR_DIR)/coverage.html tmp/main main

db-up: ## Start application database
	$(call check-env,$(ENV_FILE))
	@. $(ENV_FILE) && POSTGRES_ENV_FILE=./$(ENV_FILE) POSTGRES_PORT=$${DATABASE_PORT} docker-compose up postgres

db-down: ## Stop application database
	docker-compose stop postgres

migrate: ## Run application migrations
	$(call check-env,$(ENV_FILE))
	@. $(ENV_FILE) && MIGRATE_ENV_FILE=./$(ENV_FILE) POSTGRES_ENV_FILE=./$(ENV_FILE) docker-compose run --rm migrate

migrate-down: ## Rollback application migrations
	$(call check-env,$(ENV_FILE))
	@. $(ENV_FILE) && MIGRATE_ENV_FILE=./$(ENV_FILE) docker-compose run --rm migrate \
		-path=/migrations -database "postgres://$$DATABASE_USER:$$DATABASE_PASSWORD@postgres:5432/$$DATABASE_NAME?sslmode=disable" down

seed: migrate ## Run database seed (requires postgres: make db-up; migrate: make migrate)
	$(call check-env,$(ENV_FILE))
	@. $(ENV_FILE) && cat db/seed/populate.sql | docker-compose exec -T postgres psql -U $$DATABASE_USER -d $$DATABASE_NAME

run: deps ## Run application (no live reload)
	go run cmd/main.go

run-dev: deps ## Run application with live reload
	@command -v air > /dev/null || (echo "Error: Air not installed. Install with: go install github.com/air-verse/air@v1.64.0" && exit 1)
	air -c $(ETC_DIR)/air.toml

test: ## Run unit tests
	@. $(TEST_ENV_FILE) && go test -v -count=1 -tags=test ./internal/transport/... ./internal/usecase/... ./internal/repository/... ./internal/entity/...

coverage: ## Generate coverage report
	@mkdir -p $(VAR_DIR)
	$(call check-env,$(TEST_ENV_FILE))
	@. $(TEST_ENV_FILE) && \
		if go test -v -count=1 --tags=test -coverprofile=$(VAR_DIR)/coverage.out \
			-coverpkg=taskmanager/internal/transport/...,taskmanager/internal/usecase/...,taskmanager/internal/repository/...,taskmanager/internal/entity/... ./internal/transport/... ./internal/usecase/... ./internal/repository/... ./internal/entity/...; then \
			go tool cover -html=$(VAR_DIR)/coverage.out -o $(VAR_DIR)/coverage.html && \
			go tool cover -func=$(VAR_DIR)/coverage.out | grep total | awk '{print "Coverage: " $$3}' && \
			echo "File: $(CURDIR)/$(VAR_DIR)/coverage.html"; \
		fi
