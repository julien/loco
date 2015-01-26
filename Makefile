SHELL := /bin/bash
COVPROFILE := coverage.out

default: test

lint:
	go fmt ./...
	go vet ./...
	golint ./...

test: lint
	-@go test -coverprofile=$(COVPROFILE)

coverage: test
	-@echo ---
	go tool cover -html=$(COVPROFILE)
