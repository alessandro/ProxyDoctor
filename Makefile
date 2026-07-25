.PHONY: build test lint vet clean setup run-cli run-server fmt

build:
	go build -o bin/proxydoctor ./cmd/cli
	go build -o bin/server ./cmd/server

test:
	go test -race -coverprofile=coverage.out ./...

test-verbose:
	go test -race -v -coverprofile=coverage.out ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .
	goimports -w -local github.com/francomano/proxydoctor .

setup:
	chmod +x setup.sh && ./setup.sh

run-cli:
	go run ./cmd/cli diagnose --url $(URL)

run-server:
	go run ./cmd/server

clean:
	rm -rf bin/ coverage.out

deps:
	go mod download
	go mod verify

tidy:
	go mod tidy
