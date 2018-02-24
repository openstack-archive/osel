# Makefile for the 'osel' project.

SOURCE?=./...

go10:
	@echo "Ensuring go 1.10 is installed"
	add-apt-repository ppa:longsleep/golang-backports
	apt-get update
	apt-get install golang-go golint

env: go10
	@echo "Checking go version:"
	@go version
	@echo "Ensuring vgo is installed."
	go get -u golang.org/x/vgo

test: 
	@echo "Running tests"
	go test $(SOURCE) -cover

fmt:
	@echo "Running fmt"
	go fmt $(SOURCE)

build:
	@echo "Running build"
	vgo build
