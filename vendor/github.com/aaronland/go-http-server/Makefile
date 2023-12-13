GOMOD=vendor

lambda-example:
	if test -f main; then rm -f main; fi
	if test -f example.zip; then rm -f example.zip; fi
	GOOS=linux go build -mod $(GOMOD) -ldflags="-s -w" -o main cmd/example/main.go
	zip example.zip main
	rm -f main
