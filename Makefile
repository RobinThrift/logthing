# aggregotor make targets
PREFIX?=$(shell pwd)

.PHONY: clean test-all build

build:
	@echo "+ $@"
	@go build .

fmt:
	@echo "+ $@"
	@gofmt -s -w .

lint:
	@echo "+ $@"
	@golint ./... | grep -v vendor | tee /dev/stderr

vet:
	@echo "+ $@"
	@go vet $(shell go list ./... | grep -v vendor)

test: fmt
	@echo "+ $@"
	@go test -v -tags $(shell go list ./... | grep -v vendor)

test-all: test lint vet

clean:
	@echo "+ $@"
	@rm -rf logthing

install:
	@echo "+ $@"
	@go install .

run:
	@go run main.go

docker:
	@echo "+ $@"
	@docker build -t robinthrift/logthing:latest .
