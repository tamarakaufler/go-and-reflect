deps:
	@go mod download
	@go mod tidy

lint:
	@golangci-lint -v run

run:
	@go run cmd/reflect/main.go

all: deps lint run

.PHONY:
	deps, lint, run