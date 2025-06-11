GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli: 	
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/extrude cmd/extrude/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/inspect cmd/inspect/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/snap cmd/snap/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/review cmd/review/main.go


wasmjs:
	GOOS=js GOARCH=wasm \
		go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -tags wasmjs \
		-o www/wasm/extrude.wasm \
		cmd/extrude-wasm/main.go
