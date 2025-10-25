#!/usr/bin/env bash

set -e

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
PROTO_DIR="$PROJECT_ROOT/proto"
GO_OUT_DIR="$PROJECT_ROOT/pkg/proto"

echo "Generating Go code from protobuf definitions..."

# Создаем директорию для сгенерированного кода
mkdir -p "$GO_OUT_DIR"

# Генерируем код для каждого protobuf файла
find "$PROTO_DIR" -name "*.proto" -print0 | while IFS= read -r -d '' proto_file; do
    echo "Generating: $proto_file"

    protoc --go_out="$GO_OUT_DIR" \
           --go-grpc_out="$GO_OUT_DIR" \
           --go_opt=paths=source_relative \
           --go-grpc_opt=paths=source_relative \
           -I="$PROTO_DIR" \
           "$proto_file"
done

echo "Protobuf code generation completed!"
