BINARY=mindtick
.DEFAULT_GOAL := default
.PHONY: tidy fmt vet build run run-bin all

default: tidy fmt vet build

tidy:
	@go mod tidy

fmt:
	@go fmt ./...

vet:
	@go vet ./...

build:
	@go build -o $(BINARY)

run:
	@go run main.go

run-bin:
	@./$(BINARY)
