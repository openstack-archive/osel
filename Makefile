# Makefile for the 'osel' project.

# Note that the installation of go and vgo is accomplished by
# tools/test-setup.sh

SOURCE?=./...

env:
	@echo "Running build"
	$(HOME)/go/bin/vgo build

test: 
	@echo "Running tests"
	$(HOME)/go/bin/vgo test $(SOURCE) -cover

fmt: 
	@echo "Running fmt"
	go fmt $(SOURCE)
