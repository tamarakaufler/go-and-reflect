deps:
	@go mod download
	@go mod tidy

lint:
	@golangci-lint -v run

run-reflect:
	@go run cmd/reflect/main.go

run-env:
	@go run cmd/env/main.go

run-json:
	@go run cmd/json/main.go

all: deps lint run-reflect run-env run-json

.PHONY:
	deps, lint, run