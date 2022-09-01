OUTPUT ?= go-grpc-example
VERSION := $(shell .build/get_version.sh)

.DEFAULT_GOAL := all

generate:
	protoc protocol.proto --go_out=plugins=grpc:.

build:
	CGO_ENABLED=0 go build -ldflags "-X main.AppVersion=$(VERSION)" -o $(OUTPUT)

test:
	go test -v ./...

get-version:
	@echo $(VERSION)

all: generate build test
