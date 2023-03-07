SHELL := /bin/bash
PID=./cmd/gophermart/gophermart.pid
APP=cmd/gophermart/gophermart
PWD := $(shell pwd)

all: fmt vet test build

.PHONY: fmt
fmt: ## Format the source code
	@echo "Formatting the source code"
	go fmt ./...

.PHONY: lint
lint: ## Lint the source code. Do not forget setup linter before: "go get -u golang.org/x/lint/golint"
	@echo "Linting the source code"
	golint ./...


.PHONY: vet
vet: ## Run staticcheck. Do not forget install it before: "go install honnef.co/go/tools/cmd/staticcheck@latest".
	@echo "Checking for code issues"
	staticcheck ./...

.PHONY: build
build: ## Build
	@echo "Cleaning old binaries"
	@rm -f "${APP}"
	@rm -f "${APP}.pid"
	@echo "Building the binaries"
	@go build -o "${APP}" cmd/gophermart/main.go

.PHONY: test
test:
	@echo "Running all tests"
	@go test -mod=mod -v ./internal/handler

.PHONY: kill
kill:
	@kill `cat ${PID}` || true

.PHONY: run
run:
	@echo "Run without params"
	@${APP}  && echo $$! > ${PID}

.PHONY: help
help: ## Show current help message
	@grep -E '^[a-zA-Z-]+:.*?## .*$$' ./Makefile | sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'
