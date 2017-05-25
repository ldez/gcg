.PHONY: all

default: binary

dependencies:
	glide install

binary:
	go build

test-unit:
	go test -v -cover -coverprofile=cover.out "github.com/ldez/gcg" ;\
    go test -v -cover -coverprofile=cover.out "github.com/ldez/gcg/core" ;\
    go test -v -cover -coverprofile=cover.out "github.com/ldez/gcg/label"
