CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/aaronland/go-string
	cp -r dsn src/github.com/aaronland/go-string
	cp -r random src/github.com/aaronland/go-string
	cp *.go src/github.com/aaronland/go-string/
	if test -d vendor; then cp -r vendor/* src/; fi

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@echo "no deps"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt dsn/*.go
	go fmt random/*.go
	go fmt *.go

bin: 	self
	rm -rf bin/*
	@GOPATH=$(GOPATH) go build -o bin/randomstr cmd/randomstr.go
