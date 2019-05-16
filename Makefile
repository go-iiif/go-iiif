vendor-deps: 
	go mod vendor

fmt:
	go fmt *.go
	go fmt aws/*.go
	go fmt cache/*.go
	go fmt cmd/*.go
	go fmt compliance/*.go
	go fmt config/*.go
	go fmt http/*.go
	go fmt image/*.go
	go fmt level/*.go
	go fmt process/*.go
	go fmt profile/*.go
	go fmt service/*.go
	go fmt source/*.go
	go fmt tile/*.go
tools: 	
	go build -o bin/iiif-server cmd/iiif-server/main.go
	go build -o bin/iiif-tile-seed cmd/iiif-tile-seed/main.go
	go build -o bin/iiif-transform cmd/iiif-transform/main.go
	go build -o bin/iiif-process cmd/iiif-process/main.go
	go build -o bin/iiif-dump-config cmd/iiif-dump-config/main.go

docker-build:
	@make docker-process-build
	@make docker-server-build

docker-cli-build:
	docker build -f Dockerfile.process -t go-iiif-process .

docker-server-build:
	docker build -f Dockerfile.server -t go-iiif-server .
