#!/usr/bin/env bash

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Configure GO env
    set GOOS=linux
    kill -9 $(lsof -t -i:8003)
elif [[ "$OSTYPE" == "darwin"* ]]; then
  # Configure GO env
  set GOOS=darwin
  lsof -t -i tcp:8003 | xargs kill
else
    echo "Unknown OS"
    exit 1
fi

go get

swag init -g cmd/global-api/main.go

# Compile the internal
go build -o sen-global-api ./cmd/global-api/main.go