CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-aws
	cp -r config src/github.com/whosonfirst/go-whosonfirst-aws/
	cp -r lambda src/github.com/whosonfirst/go-whosonfirst-aws/
	cp -r s3 src/github.com/whosonfirst/go-whosonfirst-aws/
	cp -r sqs src/github.com/whosonfirst/go-whosonfirst-aws/
	cp -r session src/github.com/whosonfirst/go-whosonfirst-aws/
	cp -r util src/github.com/whosonfirst/go-whosonfirst-aws/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:	rmdeps
	@GOPATH=$(GOPATH) go get -u "github.com/aws/aws-sdk-go"
	@GOPATH=$(GOPATH) go get -u "github.com/aaronland/go-string"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-mimetypes"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt config/*.go
	go fmt lambda/*.go
	go fmt s3/*.go
	go fmt sqs/*.go
	go fmt session/*.go
	go fmt util/*.go

bin: 	self
	@GOPATH=$(GOPATH) go build -o bin/s3 cmd/s3.go
	@GOPATH=$(GOPATH) go build -o bin/secret cmd/secret.go
