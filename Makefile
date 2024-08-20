ifeq ($(wildcard .env), .env)
include .env
export
endif

check_env_file:
	@if [ -f .env ]; then \
		echo ".env file exists"; \
	else \
		echo "Error: .env file not found"; \
		exit 1; \
	fi

env:
	cp .env.example .env

local-run: check_env_file
	go run cmd/tarantool_api/main.go

run: check_env_file
	docker compose -f deployment/docker-compose.yml up -d --remove-orphans

image-build:
	docker build -f ./build/Dockerfile -t tarantool_api:latest .

test: check_env_file
	go test ./...

deploy: image-build run