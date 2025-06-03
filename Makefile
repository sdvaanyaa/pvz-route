BINARY_NAME = app

BUILD_DIR = bin

.PHONY: all update linter build start run clean

all: run

update:
	@echo "Updating dependencies"
	@go mod tidy

linter:
	@echo "Running linters"
	@golangci-lint run ./...

build:
	@echo "Building application"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/pvz/main.go

start:
	@echo "Starting application"
	@$(BUILD_DIR)/$(BINARY_NAME)

run: update linter build start

clean:
	@echo "Cleaning up"
	@rm -rf $(BUILD_DIR)
	@go clean