.PHONY: proto
proto:
	@chmod +x scripts/generate-proto.sh
	@./scripts/generate-proto.sh

.PHONY: deps
deps:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: all
all: deps proto
