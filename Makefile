SHELL := /usr/bin/env bash

.PHONY: all fmt vet tidy build run

all: fmt vet

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

build:
	go build -o bin/voltavpn ./cmd/voltavpn

run:
	go run ./cmd/voltavpn

