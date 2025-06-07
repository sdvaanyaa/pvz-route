BINARY_NAME = app
BUILD_DIR = bin
PROTO_DIR = api/proto
PROTO_OUT = api/gen

.PHONY: all update linter build start run clean proto vendor-proto bin-deps

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

run: update linter proto build start

clean:
	@echo "Cleaning up"
	@rm -rf $(BUILD_DIR)
	@go clean

bin-deps:
	@echo "Installing protoc dependencies"
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/envoyproxy/protoc-gen-validate@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

proto: vendor-proto
	@echo "Generating proto files"
	@protoc \
		-I=$(PROTO_DIR) \
		-I=vendor.protogen \
		-I=vendor.protogen/protoc-gen-openapiv2/options \
		--go_out=$(PROTO_OUT) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT) \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(PROTO_OUT) \
		--grpc-gateway_opt=paths=source_relative \
		--validate_out="lang=go,paths=source_relative:$(PROTO_OUT)" \
		--openapiv2_out=$(PROTO_OUT) \
		$(PROTO_DIR)/pvz.proto

vendor-proto: .vendor-proto/validate .vendor-proto/google/api .vendor-proto/openapiv2

.vendor-proto/validate:
	@echo "Fetching validate.proto"
	@mkdir -p vendor.protogen/validate
	@curl -sSL https://raw.githubusercontent.com/bufbuild/protoc-gen-validate/main/validate/validate.proto -o vendor.protogen/validate/validate.proto

.vendor-proto/google/api:
	@echo "Fetching google/api proto files"
	@mkdir -p vendor.protogen/google/api
	@curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto -o vendor.protogen/google/api/annotations.proto
	@curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto -o vendor.protogen/google/api/http.proto

.vendor-proto/openapiv2:
	@echo "Fetching openapiv2 proto files"
	@mkdir -p vendor.protogen/protoc-gen-openapiv2/options
	@curl -sSL https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/main/protoc-gen-openapiv2/options/annotations.proto -o vendor.protogen/protoc-gen-openapiv2/options/annotations.proto
	@curl -sSL https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/main/protoc-gen-openapiv2/options/openapiv2.proto -o vendor.protogen/protoc-gen-openapiv2/options/openapiv2.proto