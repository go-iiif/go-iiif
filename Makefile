cli-tools: 	
	go build -mod vendor -o bin/iiif-server cmd/iiif-server/main.go
	go build -mod vendor -o bin/iiif-tile-seed cmd/iiif-tile-seed/main.go
	go build -mod vendor -o bin/iiif-transform cmd/iiif-transform/main.go
	go build -mod vendor -o bin/iiif-process cmd/iiif-process/main.go
	go build -mod vendor -o bin/iiif-dump-config cmd/iiif-dump-config/main.go

docker-build:
	docker build -f Dockerfile -t go-iiif .

lambda-handlers:
	@make lambda-process
	@make lambda-transform
	@make lambda-tileseed

lambda-process:
	if test -f main; then rm -f main; fi
	if test -f iiif-process.zip; then rm -f iiif-process.zip; fi
	GOOS=linux go build -mod vendor -o main cmd/iiif-process/main.go
	zip iiif-process.zip main
	rm -f main

lambda-transform:
	if test -f main; then rm -f main; fi
	if test -f iiif-process.zip; then rm -f iiif-transform.zip; fi
	GOOS=linux go build -mod vendor -o main cmd/iiif-transform/main.go
	zip iiif-transform.zip main
	rm -f main

lambda-tileseed:
	if test -f main; then rm -f main; fi
	if test -f iiif-tile-seed.zip; then rm -f iiif-tile-seed.zip; fi
	GOOS=linux go build -mod vendor -o main cmd/iiif-tile-seed/main.go
	zip iiif-tileseed.zip main
	rm -f main

