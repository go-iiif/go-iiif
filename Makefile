CWD=$(shell pwd)

GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/iiif-server cmd/iiif-server/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/iiif-tile-seed cmd/iiif-tile-seed/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/iiif-transform cmd/iiif-transform/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/iiif-process cmd/iiif-process/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/iiif-dump-config cmd/iiif-dump-config/main.go

debug-server:
	mkdir -p fixtures/cache
	go run cmd/iiif-server/main.go \
		-config-derivatives-cache-uri file://$(CWD)/fixtures/cache \
		-config-images-source-uri file://$(CWD)/fixtures/images \
		-example \
		-verbose

lambda:
	@make lambda-handlers

docker:
	@make docker-build

docker-build:
	docker build -f Dockerfile -t go-iiif .

lambda-handlers:
	@make lambda-process
	@make lambda-transform
	@make lambda-tileseed

lambda-process:
	if test -f bootstrap; then rm -f bootstrap; fi
	if test -f iiif-process.zip; then rm -f iiif-process.zip; fi
	GOARCH=arm64 GOOS=linux go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -tags lambda.norpc -o bootstrap cmd/iiif-process/main.go
	zip iiif-process.zip bootstrap
	rm -f bootstrap

lambda-transform:
	if test -f bootstrap; then rm -f bootstrap; fi
	if test -f iiif-process.zip; then rm -f iiif-transform.zip; fi
	GOARCH=arm64 GOOS=linux go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -tags lambda.norpc -o bootstrap cmd/iiif-transform/main.go
	zip iiif-transform.zip bootstrap
	rm -f bootstrap

lambda-tileseed:
	if test -f bootstrap; then rm -f bootstrap; fi
	if test -f iiif-tile-seed.zip; then rm -f iiif-tile-seed.zip; fi
	GOARCH=arm64 GOOS=linux go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -tags lambda.norpc -o bootstrap cmd/iiif-tile-seed/main.go
	zip iiif-tileseed.zip bootstrap
	rm -f bootstrap

bump-version:
	perl -i -p -e 's/github.com\/go-iiif\/go-iiif\/$(PREVIOUS)/github.com\/go-iiif\/go-iiif\/$(NEW)/g' go.mod
	perl -i -p -e 's/github.com\/go-iiif\/go-iiif\/$(PREVIOUS)/github.com\/go-iiif\/go-iiif\/$(NEW)/g' README.md
	find . -name '*.go' | xargs perl -i -p -e 's/github.com\/go-iiif\/go-iiif\/$(PREVIOUS)/github.com\/go-iiif\/go-iiif\/$(NEW)/g'
