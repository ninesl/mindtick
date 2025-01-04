BINARY=mindtick
.DEFAULT_GOAL := all
.PHONY: fmt vet build run run-bin all

all: fmt vet build run-bin

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
