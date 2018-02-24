# Makefile for the 'osel' project.

SOURCE?=./...

version:
	@echo "Checking go version:"
	@go version

test: 
	@echo "Running tests"
	go test $(SOURCE) -cover

fmt:
	@echo "Running fmt"
	go fmt $(SOURCE)

env: version fmt test
	@echo "Running env"
