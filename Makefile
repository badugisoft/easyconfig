#! /usr/bin/make

.PHONY: test

test:
	@go get
	@go test -v

bindata:
	@go-bindata -o test/bindata.go -pkg test -ignore ".+\.go" test

deps:
	@go get -u github.com/jteeuwen/go-bindata/...
