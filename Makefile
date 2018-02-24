# Makefile for the 'osel' project.

SOURCE?=./...

go10:
	@echo "Ensuring go 1.10 is installed"
	sudo add-apt-repository ppa:longsleep/golang-backports
	sudo apt-get update
	sudo apt-get install golang-go golint

env: go10
	@echo "Checking go version:"
	@go version
	@echo "Ensuring vgo is installed."
	go get -u golang.org/x/vgo

test: build 
	@echo "Running tests"
	go test $(SOURCE) -cover

fmt:
	@echo "Running fmt"
	go fmt $(SOURCE)

build:
	@echo "Running build"
	vgo build
