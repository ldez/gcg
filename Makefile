.PHONY: all

default: binary

dependencies:
	glide install

binary:
	go build