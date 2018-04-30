PROJECT_NAME := "lyra"
PKG := "github.com/azohra/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GOPATH := $(shell go env GOPATH)

.PHONY: test race msan coverage dep localbuild buildbins clean install binaries help build

test: ## Run all unit tests verbosely
	@go test -v ./...

race: ## Run data race detector
	@go test -race -short ${PKG_LIST}

msan: ## Run memory sanitizer
	@go test -msan -short ${PKG_LIST}

coverage: ## Generate global code coverage report
	@go test --coverprofile -r ./...

dep: ## Get the dependencies
	@dep ensure -v

localbuild:  # Build the binary file
	@go install ./cmd/${PROJECT_NAME}/...
	@ln -s  ${GOPATH}/bin/${PROJECT_NAME} /usr/local/bin

buildbins: # Build on diff platform
	@env GOOS=linux GOARCH=amd64 go build -v -o build/bin/linux/amd64/${PROJECT_NAME} ./cmd/${PROJECT_NAME}/...
	@env GOOS=darwin GOARCH=amd64 go build -v -o build/bin/darwin/amd64/${PROJECT_NAME} ./cmd/${PROJECT_NAME}/...
	@env GOOS=windows GOARCH=amd64 go build -v -o build/bin/win/amd64/${PROJECT_NAME}.exe ./cmd/${PROJECT_NAME}/...

clean: ## Remove previous build and undo install
	@rm -f /usr/local/bin/${PROJECT_NAME}
	@rm -f ${GOPATH}/bin/${PROJECT_NAME}
	@rm -rf build/bin
	@rm -f ./lyra

build: ## Build a binary in the current working directory
	@go build ./cmd/lyra/...

install: dep test localbuild ## Build app and install it into shell PATH

binaries: dep test buildbins # Build diff binaries

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'