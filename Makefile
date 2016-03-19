default: all

all: test lint vet

test: go test ./...

fmt: gofmt -s -d ./...

lint: golint ./...

vet: go vet ./...

loc:
	wc -l *.go
