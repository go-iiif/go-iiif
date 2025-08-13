GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/aws-sts-session cmd/aws-sts-session/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/aws-mfa-session cmd/aws-mfa-session/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/aws-get-credentials cmd/aws-get-credentials/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/aws-cognito-credentials cmd/aws-cognito-credentials/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/aws-set-env cmd/aws-set-env/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/aws-sign-request cmd/aws-sign-request/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/aws-credentials-json-to-ini cmd/aws-credentials-json-to-ini/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/aws-imds-credentials cmd/aws-imds-credentials/main.go
