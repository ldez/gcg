.PHONY: clean test-unit build check fmt

GOFILES := $(shell go list -f '{{range $$index, $$element := .GoFiles}}{{$$.Dir}}/{{$$element}}{{"\n"}}{{end}}' ./... | grep -v '/vendor/')

default: clean check test-unit build

dependencies:
	dep ensure

build:
	go build

test-unit:
	go test -v ./...

check:
	golangci-lint run

fmt:
	gofmt -s -l -w $(GOFILES)

clean:
	rm -f gcg cover.out