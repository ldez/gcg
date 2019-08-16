.PHONY: clean test build check fmt

export GO111MODULE=on

GOFILES := $(shell go list -f '{{range $$index, $$element := .GoFiles}}{{$$.Dir}}/{{$$element}}{{"\n"}}{{end}}' ./... | grep -v '/vendor/')

default: clean check test build

build:
	go build

test: clean
	go test -v ./...

check:
	golangci-lint run

fmt:
	gofmt -s -l -w $(GOFILES)

clean:
	rm -f gcg cover.out