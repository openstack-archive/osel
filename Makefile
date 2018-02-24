# Makefile for the 'osel' project.

SOURCE?=./...

env:
	@echo "Checking go version:"
	@go version

test: 
	@echo "Running tests"
	go test $(SOURCE) -cover

fmt:
	@echo "Running fmt"
	go fmt $(SOURCE)

