.PHONY: all

default: clean test-unit validate build

dependencies:
	dep ensure

build:
	go build

validate:
	./.script/make.sh validate-gofmt validate-govet validate-golint validate-misspell

test-unit:
	./.script/make.sh test-unit

clean:
	rm -f gcg cover.out