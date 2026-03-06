GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

GOROOT=$(shell go env GOROOT)

tag-data:
	curl -o tags/tags_data.go 

rebuild-wasm:
	@make wasm

wasm:
	GOOS=js GOARCH=wasm go build -ldflags="$(LDFLAGS)" -mod $(GOMOD) -tags wasm -o www/wasm/update_exif.wasm cmd/update-exif-wasm/main.go
	GOOS=js GOARCH=wasm go build -ldflags="$(LDFLAGS)" -mod $(GOMOD) -tags wasm -o www/wasm/supported_tags.wasm cmd/tags-supported-wasm/main.go

cli:
	@make wasm
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/tags-is-supported cmd/tags-is-supported/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/tags-supported cmd/tags-supported/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/update-exif cmd/update-exif/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/update-exif-wasm-demo cmd/update-exif-wasm-demo/main.go

debug:
	go run -mod $(GOMOD) cmd/update-exif-wasm-demo/main.go

lambda:
	@make lambda-server

lambda-server:
	if test -f main; then rm -f main; fi
	if test -f server.zip; then rm -f server.zip; fi
	GOOS=linux go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o main cmd/update-exif-server/main.go
	zip server.zip main
	rm -f main
