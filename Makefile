#!/usr/bin/make -f

SHELL := /bin/bash

BINARY_NAME := simple-config-server
BUILD_DIR := bin
CONFIG_DIR := configurations
PORT ?= 8080
JWT_SECRET ?= secret
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

OS_LIST := linux darwin windows
ARCH_LIST := amd64 arm64

all: build

build-all: clean
	@echo "🔧 Building binaries for all supported platforms..."
	@mkdir -p $(BUILD_DIR)
	@for os in $(OS_LIST); do \
		for arch in $(ARCH_LIST); do \
			echo "🚀 Building for $$os/$$arch..."; \
			GOOS=$$os GOARCH=$$arch go build -v -o $(BUILD_DIR)/$(BINARY_NAME)-$$os-$$arch -ldflags="-w -s" .; \
		done; \
	done
	@echo "✅ All builds complete."

build: 
	@echo "🔨 Building the application..."
	@mkdir -p $(BUILD_DIR)
	@go build -v -o $(BUILD_DIR)/$(BINARY_NAME) -ldflags="-w -s" .
	@echo "✅ Build complete."

run: build
	@echo "🚀 Running the application on port $(PORT)..."
	@PORT=$(PORT) JWT_SECRET=$(JWT_SECRET) ./$(BUILD_DIR)/$(BINARY_NAME)

dev: build
	@echo "🚀 Running the application on port $(PORT)..."
	@cd $(BUILD_DIR) && PORT=$(PORT) JWT_SECRET=$(JWT_SECRET) ./$(BINARY_NAME)

clean:
	@echo "🧹 Cleaning the build directory..."
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete."

fmt:
	@echo "🖋 Formatting the code..."
	@go fmt ./...
	@echo "✅ Format complete."

vet:
	@echo "🔍 Running go vet..."
	@go vet ./...
	@echo "✅ Vet complete."

lint: install_deps
	@echo "🔍 Linting the code..."
	@golangci-lint run || { echo "❌ Linting failed!"; exit 1; }
	@echo "✅ Lint complete."

install_deps:
	@echo "📦 Installing dependencies..."
	@command -v golangci-lint >/dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Dependencies installed."

help:
	@echo "📜 Available Makefile commands:"
	@echo "  all        - Build the application (default)"
	@echo "  build      - Compile the Go application"
	@echo "  run        - Run the application"
	@echo "  clean      - Remove built binaries"
	@echo "  fmt        - Format the Go source files"
	@echo "  vet        - Analyze the code with go vet"
	@echo "  lint       - Run static analysis with golangci-lint"
	@echo "  test       - Run unit tests with coverage"
	@echo "  install_deps - Install necessary dependencies (golangci-lint)"
	@echo "  help       - Display this help message"

.PHONY: all build run clean fmt vet lint test install_deps help
