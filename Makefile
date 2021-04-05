deps:
	@go mod download
	@go mod tidy

lint:
	@golangci-lint -v run

run-reflect:
	@go run cmd/reflect/main.go

run-env:
	@go run cmd/env/main.go

run-marshal:
	@go run cmd/marshal/main.go

all: deps lint run-reflect run-env run-marshal

.PHONY:
	deps, lint, run