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

run: check_env_file
	go run cmd/tarantool_api/main.go

test: check_env_file
	go test ./...