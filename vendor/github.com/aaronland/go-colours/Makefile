CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/aaronland/go-colours
	cp *.go src/github.com/aaronland/go-colours/
	cp -r grid src/github.com/aaronland/go-colours/
	cp -r extruder src/github.com/aaronland/go-colours/
	cp -r palette src/github.com/aaronland/go-colours/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/nfnt/resize"
	@GOPATH=$(GOPATH) go get -u "github.com/neocortical/noborders"
	@GOPATH=$(GOPATH) go get -u "github.com/lucasb-eyer/go-colorful"
	@GOPATH=$(GOPATH) go get -u "github.com/RobCherry/vibrant"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt *.go
	go fmt cmd/*.go
	go fmt grid/*.go
	go fmt extruder/*.go
	go fmt palette/*.go

bin: 	self
	rm -rf bin/*
	@GOPATH=$(GOPATH) go build -o bin/extrude cmd/extrude.go
	@GOPATH=$(GOPATH) go build -o bin/inspect cmd/inspect.go
	@GOPATH=$(GOPATH) go build -o bin/snap cmd/snap.go

