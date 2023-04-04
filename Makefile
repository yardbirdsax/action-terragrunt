GOBIN_PATH = $$PWD/.bin
ENV_VARS = GOBIN="$(GOBIN_PATH)" PATH="$(GOBIN_PATH):$$PATH"
convey:
	$(ENV_VARS) goconvey
coverage:
	go tool cover -func=coverage.out
generate:
	$(ENV_VARS) go generate ./...
test:
	go test -count=1 -coverprofile coverage.out ./...
tidy:
	go mod tidy
tools:
	$(ENV_VARS) go install $$(go list -f '{{join .Imports " "}}' tools.go)
build:
	CGO_ENABLED=0 go build -o dist/app main.go
template:
	if [ -z "$$GIT_TAG" ]; then \
		GIT_TAG=$$(git rev-parse HEAD); \
		export GIT_TAG; \
	fi && \
	docker run -e GIT_TAG -e GITHUB_OWNER -v $$PWD:/workspace hairyhenderson/gomplate:stable -f /workspace/action.yaml.tpl -o /workspace/action.yaml
