.PHONY: clean test build check

export GO111MODULE=on

default: clean check test build

build:
	go build

test: clean
	go test -v ./...

check:
	golangci-lint run

clean:
	rm -f gcg cover.out