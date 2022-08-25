OUTPUT ?= go-grpc-example

.DEFAULT_GOAL := all

generate:
	protoc protocol.proto --go_out=plugins=grpc:.

build:
	CGO_ENABLED=0 go build -o $(OUTPUT)

test:
	go test -v ./...

all: build test
