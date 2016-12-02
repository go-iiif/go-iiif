CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src/github.com/thisisaaronland/go-iiif; then rm -rf src/github.com/thisisaaronland/go-iiif; fi
	mkdir -p src/github.com/thisisaaronland/go-iiif
	cp iiif.go src/github.com/thisisaaronland/go-iiif/
	cp -r aws src/github.com/thisisaaronland/go-iiif/
	cp -r cache src/github.com/thisisaaronland/go-iiif/
	cp -r compliance src/github.com/thisisaaronland/go-iiif/
	cp -r config src/github.com/thisisaaronland/go-iiif/
	cp -r image src/github.com/thisisaaronland/go-iiif/
	cp -r level src/github.com/thisisaaronland/go-iiif/
	cp -r profile src/github.com/thisisaaronland/go-iiif/
	cp -r source src/github.com/thisisaaronland/go-iiif/
	cp -r tile src/github.com/thisisaaronland/go-iiif/
	cp -r vendor/src/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:	rmdeps
	@GOPATH=$(GOPATH) go get -u "github.com/facebookgo/grace/gracehttp"
	@GOPATH=$(GOPATH) go get -u "github.com/fogleman/primitive"
	@GOPATH=$(GOPATH) go get -u "github.com/gorilla/context"
	@GOPATH=$(GOPATH) go get -u "github.com/gorilla/mux"
	@GOPATH=$(GOPATH) go get -u "github.com/patrickmn/go-cache"
	@GOPATH=$(GOPATH) go get -u "golang.org/x/image/tiff"
	@GOPATH=$(GOPATH) go get -u "golang.org/x/image/webp"
	@GOPATH=$(GOPATH) go get -u "gopkg.in/h2non/bimg.v1"
	@GOPATH=$(GOPATH) go get -u "github.com/koyachi/go-atkinson"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-sanitize"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-csv"
	@GOPATH=$(GOPATH) go get -u "github.com/jtacoma/uritemplates"
	@GOPATH=$(GOPATH) go get -u "github.com/aws/aws-sdk-go"

vendor-deps: deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor/src; then rm -rf vendor/src; fi
	cp -r src vendor/src
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt *.go
	go fmt aws/*.go
	go fmt cache/*.go
	go fmt cmd/*.go
	go fmt compliance/*.go
	go fmt image/*.go
	go fmt level/*.go
	go fmt profile/*.go
	go fmt source/*.go
	go fmt tile/*.go

bin: 	self
	@GOPATH=$(GOPATH) go build -o bin/iiif-server cmd/iiif-server.go
	@GOPATH=$(GOPATH) go build -o bin/iiif-tile-seed cmd/iiif-tile-seed.go
	@GOPATH=$(GOPATH) go build -o bin/iiif-transform cmd/iiif-transform.go
	@GOPATH=$(GOPATH) go build -o bin/iiif-dump-config cmd/iiif-dump-config.go
