cli-tools: 	
	go build -mod vendor -o bin/iiif-server cmd/iiif-server/main.go
	go build -mod vendor -o bin/iiif-tile-seed cmd/iiif-tile-seed/main.go
	go build -mod vendor -o bin/iiif-transform cmd/iiif-transform/main.go
	go build -mod vendor -o bin/iiif-process cmd/iiif-process/main.go
	go build -mod vendor -o bin/iiif-dump-config cmd/iiif-dump-config/main.go

docker-build:
	@make docker-process-build
	@make docker-server-build

docker-process-build:
	docker build -f Dockerfile.process -t go-iiif-process .

docker-server-build:
	docker build -f Dockerfile.server -t go-iiif-server .

lambda-handlers:
	@make lambda-process

lambda-process:
	if test -f main; then rm -f main; fi
	if test -f iiif-process.zip; then rm -f iiif-process.zip; fi
	GOOS=linux go build -mod vendor -o main cmd/iiif-process/main.go
	zip iiif-process.zip main
	rm -f main

