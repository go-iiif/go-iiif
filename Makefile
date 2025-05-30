CWD=$(shell pwd)

GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/iiif-server cmd/iiif-server/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/iiif-tile-seed cmd/iiif-tile-seed/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/iiif-transform cmd/iiif-transform/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/iiif-process cmd/iiif-process/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/iiif-dump-config cmd/iiif-dump-config/main.go

debug-seed:
	if test -d $(CWD)/fixtures/cache/spank; then rm -rf $(CWD)/fixtures/cache/spank; fi
	go run -mod $(GOMOD) cmd/iiif-tile-seed/main.go \
		-config-images-source-uri file://$(CWD)/fixtures/images \
		-config-derivatives-cache-uri file://$(CWD)/fixtures/cache \
		-verbose \
		-generate-html \
		'rewrite:///spanking-cat.jpg?target=spank'

debug-seed-csv:
	if test -d $(CWD)/fixtures/cache/spanking-csv; then rm -rf $(CWD)/fixtures/cache/spanking-csv; fi
	if test -d $(CWD)/fixtures/cache/walrus-csv; then rm -rf $(CWD)/fixtures/cache/walrus-csv; fi
	if test -f $(CWD)/fixtures/seed.csv; then $(CWD)/fixtures/seed.csv; fi
	echo "source_filename,source_root,target_filename,target_root" > $(CWD)/fixtures/seed.csv
	echo "spanking-cat.jpg,$(CWD)/fixtures/images,spanking-csv,$(CWD)/fixtures/cache" >> $(CWD)/fixtures/seed.csv
	echo "walrus.jpg,$(CWD)/fixtures/images,walrus-csv,$(CWD)/fixtures/cache" >> $(CWD)/fixtures/seed.csv
	go run -mod $(GOMOD) cmd/iiif-tile-seed/main.go \
		-mode csv \
		-generate-html \
		-verbose \
		$(CWD)/fixtures/seed.csv

debug-process:
	if test -d $(CWD)/fixtures/cache/999; then rm -rf $(CWD)/fixtures/cache/999; fi
	go run -mod $(GOMOD) cmd/iiif-process/main.go \
		-config-derivatives-cache-uri file://$(CWD)/fixtures/cache \
		-config-images-source-uri file://$(CWD)/fixtures/images \
		-report \
		-report-bucket-uri file://$(CWD)/fixtures/reports \
		-report-html \
		-verbose \
		'idsecret:///spanking-cat.jpg?id=9998&secret=abc&secret_o=def&format=jpg&label=x'

debug-server:
	mkdir -p fixtures/cache
	go run -mod $(GOMOD) cmd/iiif-server/main.go \
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
