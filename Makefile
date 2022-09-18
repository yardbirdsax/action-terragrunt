GOBIN_PATH = $$PWD/.bin
ENV_VARS = GOBIN="$(GOBIN_PATH)" PATH="$(GOBIN_PATH):$$PATH"
convey:
	$(ENV_VARS) goconvey
coverage:
	go tool cover -func=coverage.out
generate:
	$(ENV_VARS) go generate ./...
test:
	go test -v -count=1 -coverprofile coverage.out ./...
tidy:
	go mod tidy
tools:
	$(ENV_VARS) go install $$(go list -f '{{join .Imports " "}}' tools.go)
build: test
	CGO_ENABLED=0 go build -o app main.go
