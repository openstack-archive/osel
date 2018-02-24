# Makefile for the 'osel' project.

SOURCE?=./...

go10:
	@echo "Ensuring go 1.10 is installed"
	sudo add-apt-repository ppa:longsleep/golang-backports
	sudo apt-get update
	sudo apt-get install -y golang-go golint

env10: go10
	@echo "Checking go version:"
	@go version
	@echo "Ensuring vgo is installed."
	go get -u golang.org/x/vgo

test10: build10
	@echo "Running tests"
	go test $(SOURCE) -cover

fmt10: env10
	@echo "Running fmt"
	go fmt $(SOURCE)

build10: env10
	@echo "Running build"
	vgo build


env:
	@/bin/true

fmt:
	@/bin/true

test:
	@/bin/true
