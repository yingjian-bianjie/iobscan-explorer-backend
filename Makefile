#!/usr/bin/make -f

build: go.sum
ifeq ($(OS),Windows_NT)
	go build  -o build/ddcparser.exe ./cmd/ddcparser
else
	go build  -o build/ddcparser ./cmd/ddcparser
endif

build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

install:
	go install ./cmd/ddcparser

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs misspell -w

