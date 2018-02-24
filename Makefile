# Makefile for the 'osel' project.

# Note that the installation of go and vgo is accomplished by
# tools/test-setup.sh

SOURCE?=./...

env:
	@echo "Running build"
	vgo build

test: 
	@echo "Running tests"
	go test $(SOURCE) -cover

fmt: 
	@echo "Running fmt"
	go fmt $(SOURCE)
