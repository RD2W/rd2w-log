.DEFAULT_GOAL := help

.PHONY: proto
proto: check-deps
	@echo "🚀 Launch script for Protocol Buffer code generation..."
	@chmod +x scripts/generate-proto.sh
	@./scripts/generate-proto.sh

.PHONY: deps
deps:
	@echo "🛠️ Installing protobuf dependencies..."
	@which protoc > /dev/null || (echo "⚠️  Note: protoc not found. Please install protobuf-compiler" && sleep 2)
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: check-deps
check-deps:
	@echo "🔍 Checking build dependencies..."
	@which protoc > /dev/null || (echo "❌ Error: protoc not installed.\n   On Ubuntu: sudo apt-get install protobuf-compiler\n   On macOS: brew install protobuf" && exit 1)
	@[ -f "$(shell go env GOPATH)/bin/protoc-gen-go" ] || (echo "❌ Error: protoc-gen-go not installed. Run 'make deps'" && exit 1)
	@[ -f "$(shell go env GOPATH)/bin/protoc-gen-go-grpc" ] || (echo "❌ Error: protoc-gen-go-grpc not installed. Run 'make deps'" && exit 1)
	@echo "✅ All dependencies are satisfied"

.PHONY: all
all: deps proto
	@echo "✅ Build setup completed!"

.PHONY: clean-proto
clean-proto:
	@echo "🧹 Cleaning generated protobuf code..."
	@rm -rf pkg/proto/*

.PHONY: test
test:
	@go test ./...

.PHONY: build
build: proto
	@go build ./...

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  deps       - Install Go protobuf dependencies (requires protoc)"
	@echo "  proto      - Generate protobuf code"
	@echo "  all        - Install deps and generate proto"
	@echo "  test       - Run tests"
	@echo "  clean-proto - Remove generated protobuf code"
	@echo ""
	@echo "⚠️ Note: protoc must be installed separately via system package manager"
