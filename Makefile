SHELL := bash

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint: tidy
	go vet ./...
	go fmt ./...

.PHONY: unit-test
unit-test:
	go test ./...
