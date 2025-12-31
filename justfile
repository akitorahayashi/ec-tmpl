# ==============================================================================
# justfile for ec-tmpl automation
# ==============================================================================

set dotenv-load
set shell := ["bash", "-cu"]

APP_NAME := env("EC_TMPL_APP_NAME", "ec-tmpl")
HOST_IP := env("EC_TMPL_DEV_HOST_IP", "127.0.0.1")
DEV_PORT := env("EC_TMPL_DEV_PORT", "8000")
HOST_PORT := env("EC_TMPL_HOST_PORT", "8000")

DEV_COMPOSE_PROJECT := env("EC_TMPL_DEV_PROJECT_NAME", "ec-tmpl-dev")
PROD_COMPOSE_PROJECT := env("EC_TMPL_PROD_PROJECT_NAME", "ec-tmpl-prod")

# default target
default: help

# Show available recipes
help:
	@echo "Usage: just [recipe]"
	@echo "Available recipes:"
	@just --list | tail -n +2 | awk '{printf "  \033[36m%-20s\033[0m %s\n", $1, substr($0, index($0, $2))}'

# ==============================================================================
# Environment Setup
# ==============================================================================

setup: # Initialize local environment (tools + .env)
	@if [ ! -f .env ] && [ -f .env.example ]; then cp .env.example .env; fi
	@mkdir -p .bin
	@GOBIN="$(pwd)/.bin" go install golang.org/x/tools/cmd/goimports@v0.30.0
	@GOBIN="$(pwd)/.bin" go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

# ==============================================================================
# Development Environment Commands
# ==============================================================================

dev: # Local development server (in-process, no Docker)
	@if [ ! -f .bin/air ]; then GOBIN="$(pwd)/.bin" go install github.com/air-verse/air@v1.52.3; fi
	@EC_TMPL_BIND_IP="{{HOST_IP}}" EC_TMPL_BIND_PORT="{{DEV_PORT}}" .bin/air -c .air.toml

up: # Development stack (docker compose)
	@COMPOSE_PROJECT_NAME="{{DEV_COMPOSE_PROJECT}}" EC_TMPL_HOST_PORT="{{HOST_PORT}}" docker compose up -d api

down: # Stop development stack (docker compose)
	@COMPOSE_PROJECT_NAME="{{DEV_COMPOSE_PROJECT}}" docker compose down --remove-orphans

rebuild: # Rebuild development stack (docker compose)
	@COMPOSE_PROJECT_NAME="{{DEV_COMPOSE_PROJECT}}" docker compose down --remove-orphans
	@COMPOSE_PROJECT_NAME="{{DEV_COMPOSE_PROJECT}}" docker compose build --no-cache api
	@COMPOSE_PROJECT_NAME="{{DEV_COMPOSE_PROJECT}}" docker compose up -d api

up-prod: # Production stack (docker compose)
	@COMPOSE_PROJECT_NAME="{{PROD_COMPOSE_PROJECT}}" EC_TMPL_HOST_PORT="{{HOST_PORT}}" docker compose up -d api-prod

down-prod: # Stop production stack (docker compose)
	@COMPOSE_PROJECT_NAME="{{PROD_COMPOSE_PROJECT}}" docker compose down --remove-orphans

rebuild-prod: # Rebuild production stack (docker compose)
	@COMPOSE_PROJECT_NAME="{{PROD_COMPOSE_PROJECT}}" docker compose down --remove-orphans
	@COMPOSE_PROJECT_NAME="{{PROD_COMPOSE_PROJECT}}" docker compose build --no-cache api-prod
	@COMPOSE_PROJECT_NAME="{{PROD_COMPOSE_PROJECT}}" docker compose up -d api-prod

# ==============================================================================
# CODE QUALITY
# ==============================================================================

# Automatically format and fix code
fix:
	@echo "🔧 Formatting and fixing code..."
	@goimports -w .
	@golangci-lint run --fix

# Run static checks (format, lint)
check: fix
	@echo "🔍 Running static checks..."
	@test -z "$(goimports -l .)" || (echo "goimports check failed" && exit 1)
	@golangci-lint run

# ==============================================================================
# TESTING
# ==============================================================================

# Run complete test suite
test:
	@just local-test
	@just docker-test
	@echo "✅ All tests passed!"

# Run lightweight local (in-process) test suite
local-test:
	@just unit-test
	@just intg-test
	@echo "✅ All local tests passed!"

# Run unit tests
unit-test:
	@echo "🚀 Running unit tests..."
	@go test ./...

# Run integration tests (in-process)
intg-test:
	@echo "🚀 Running integration tests..."
	@go test -tags=intg ./...

# Run all Docker-based tests
docker-test:
	@just api-test
	@just e2e-test
	@echo "✅ All Docker tests passed!"

# Run dockerized API tests (development target)
api-test:
	@echo "🚀 Building image for dockerized API tests (development target)..."
	@docker build --target development -t ec-tmpl-test:dev .
	@echo "🚀 Running dockerized API tests (development target)..."
	@EC_TMPL_E2E_IMAGE=ec-tmpl-test:dev go test -tags=api ./tests/api

# Run e2e tests (production target)
e2e-test:
	@echo "🚀 Building image for production acceptance tests..."
	@docker build --target production -t ec-tmpl-test:prod .
	@echo "🚀 Running production acceptance tests..."
	@EC_TMPL_E2E_IMAGE=ec-tmpl-test:prod go test -tags=e2e ./tests/e2e

# ==============================================================================
# CLEANUP
# ==============================================================================

clean: # Remove generated artifacts
	@rm -rf .bin tmp .tmp coverage.out
