prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-log; then rm -rf src/github.com/whosonfirst/go-whosonfirst-log; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-log
	cp *.go src/github.com/whosonfirst/go-whosonfirst-log/

deps:   self

fmt:
	go fmt *.go

test:	self fmt
	@GOPATH=$(shell pwd) go build -o bin/test cmd/test.go

bin: 	test

