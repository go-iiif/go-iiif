cli:
	go build -mod vendor -o bin/lookup cmd/lookup/main.go

apache-lookup-table:
	curl -s -o tables/apache.txt 'https://svn.apache.org/viewvc/httpd/httpd/branches/2.4.x/docs/conf/mime.types?view=co'

lookup-tables:
	if test -d lookup; then rm -rf lookup; fi
	mkdir lookup
	go run -mod vendor cmd/build-lookup-tables/main.go -lookup extension tables/apache.txt tables/custom.txt > lookup/extension.go
	go run -mod vendor cmd/build-lookup-tables/main.go -lookup mimetype tables/apache.txt tables/custom.txt > lookup/mimetype.go

rebuild:
	@make apache-lookup-table
	@make lookup-tables
