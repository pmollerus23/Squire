BINARY_NAME=agent-cli
VERSION=1.0.0

.PHONY: build
build:
	go build -ldflags="-s -w" -o bin/$(BINARY_NAME) .

.PHONY: build-all
build-all:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(BINARY_NAME)-windows-amd64.exe .

.PHONY: install
install:
	go install -ldflags="-s -w" .

.PHONY: clean
clean:
	rm -rf bin/

.PHONY: run
run:
	go run .