vendor-deps:
	go mod vendor

fmt:
	go fmt *.go
	go fmt cmd/*.go
	go fmt grid/*.go
	go fmt extruder/*.go
	go fmt palette/*.go

tools: 	
	go build -o bin/extrude cmd/extrude/main.go
	go build -o bin/inspect cmd/inspect/main.go
	go build -o bin/snap cmd/snap/main.go

