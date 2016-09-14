prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-csv; then rm -rf src/github.com/whosonfirst/go-whosonfirst-csv; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-csv
	cp csv.go src/github.com/whosonfirst/go-whosonfirst-csv/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	rmdeps deps fmt bin

deps:   self

fmt:
	go fmt *.go
	go fmt cmd/*.go

bin: 	self
	@GOPATH=$(shell pwd) go build -o bin/wof-csv-filter cmd/wof-csv-filter.go
