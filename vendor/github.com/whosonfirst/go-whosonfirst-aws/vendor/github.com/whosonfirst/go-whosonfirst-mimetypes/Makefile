CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-mimetypes
	cp *.go src/github.com/whosonfirst/go-whosonfirst-mimetypes/
	cp -r lookup src/github.com/whosonfirst/go-whosonfirst-mimetypes/
	if test -d vendor; then cp -r vendor/* src/; fi

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

docker-build:
	docker build -t wof-readwrited .

deps:
	@echo "no dependencies yet"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt lookup/*.go
	go fmt *.go

bin: 	self
	@GOPATH=$(GOPATH) go build -o bin/wof-mimetype-lookup cmd/wof-mimetype-lookup.go

apache-lookup-table:
	curl -s -o tables/apache.txt 'https://svn.apache.org/viewvc/httpd/httpd/branches/2.2.x/docs/conf/mime.types?view=co'

lookup-tables:	self
	@GOPATH=$(GOPATH) go build -o bin/wof-mimetype-build-lookup-tables cmd/wof-mimetype-build-lookup-tables.go
	if test -d lookup; then rm -rf lookup; fi
	mkdir lookup
	bin/wof-mimetype-build-lookup-tables -lookup extension tables/apache.txt tables/custom.txt > lookup/extension.go
	bin/wof-mimetype-build-lookup-tables -lookup mimetype tables/apache.txt tables/custom.txt > lookup/mimetype.go
	go fmt lookup/*.go
	rm bin/wof-mimetype-build-lookup-tables

rebuild:
	@make apache-lookup-table
	@make lookup-tables
	@make bin
