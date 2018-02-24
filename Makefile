# Makefile for the 'osel' project.

SOURCE?=./...

test: 
	go test $(SOURCE) -cover

fmt:
	go fmt $(SOURCE)

env: test fmt
