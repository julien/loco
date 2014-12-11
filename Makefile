SHELL := /bin/bash
COVPROFILE := coverage.out

default: test

lint:
	-@echo --- Linting
	-@go vet ./...
	-@golint ./...

test: lint
	-@echo --- Launcing tests
	-@go test -coverprofile=$(COVPROFILE)

coverage: test
	go tool cover -html=$(COVPROFILE)

