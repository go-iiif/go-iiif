docker-build:
	docker build -t wof-readwrited .

deps:
	go mod vendor

build:
	go build ./...

fmt:
	go fmt cmd/*.go
	go fmt lookup/*.go
	go fmt *.go

tools: 	
	go build -o bin/wof-mimetype-lookup cmd/wof-mimetype-lookup/main.go

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
	@make tools
