#!/bin/bash

ENV=${1:-dev}

templ generate

WASM_OUTPUT="./public/lib.wasm"
WASM_INPUT="./front/web.go"

if [ "$ENV" = "prod" ]; then
    echo "Building for production..."
    GOOS=js GOARCH=wasm go build -tags=prod -ldflags="-s -w" -o $WASM_OUTPUT $WASM_INPUT
else
    echo "Building for development..."
    GOOS=js GOARCH=wasm go build -tags=dev -o $WASM_OUTPUT $WASM_INPUT
fi

tailwindcss -i ./global.css -o ./public/output.css

go run ./server/server.go
