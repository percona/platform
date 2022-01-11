default: help

help:                                      ## Display this help message
	@echo "Please use \`make <target>\` where <target> is one of:"
	@grep '^[a-zA-Z]' $(MAKEFILE_LIST) | \
		awk -F ':.*?## ' 'NF==2 {printf "  %-26s%s\n", $$1, $$2}'

init:                                      ## Install development tools
	cd tools && go generate -x -tags=tools

ci-init:                                   ## Initialize CI environment

gen:                                       ## Format, check, and generate code using buf; TODO Add lint and break commands
	rm -rf gen
	bin/buf generate
	make format

gen-dev:                                   ## Keep it to make the CI green, TODO remove it in the next PR
	make gen

gen-code:                                  ## Generate code
	go generate ./...
	go install ./...

format:                                    ## Format source code
	bin/gofumpt -l -w .
	bin/goimports -local github.com/percona-platform/platform -l -w .

check:                                     ## Run checks/linters for the whole project
	bin/go-consistent -exclude=tools -pedantic ./...
	bin/golangci-lint run

test:                                      ## Run tests
	go test -race ./...

test-cover:                                ## Run tests and collect per-package coverage information
	go test -race -timeout=10m -count=1 -coverprofile=cover.out -covermode=atomic ./...

test-crosscover:                           ## Run tests and collect cross-package coverage information
	go test -race -timeout=10m -count=1 -coverprofile=crosscover.out -covermode=atomic -p=1 -coverpkg=./... ./...

descriptors:                               ## Update files used for breaking changes detection
	bin/buf build -o platform.bin --as-file-descriptor-set

saas:                                      ## Extract public APIs and generated files into ../saas
	go run post-processing.go -project saas

fuzz-check-build:
	bin/go-fuzz-build -o pkg/check/check-fuzz.zip github.com/percona-platform/platform/pkg/check

fuzz-check-data: fuzz-check-build          ## Fuzz data tests
	bin/go-fuzz -workdir pkg/check/fuzzdata -bin pkg/check/check-fuzz.zip -func FuzzData

fuzz-check-signature: fuzz-check-build     ## Fuzz signature tests
	bin/go-fuzz -workdir pkg/check/fuzzdata -bin pkg/check/check-fuzz.zip -func FuzzSign

fuzz-check-pubkey: fuzz-check-build        ## Fuzz public key tests
	bin/go-fuzz -workdir pkg/check/fuzzdata -bin pkg/check/check-fuzz.zip -func FuzzPublicKey

fuzz-starlark:                             ## Fuzz starlark package
	go test -count=1 github.com/percona-platform/platform/pkg/starlark
	bin/go-fuzz-build -o pkg/starlark/starlark-fuzz.zip github.com/percona-platform/platform/pkg/starlark
	bin/go-fuzz -workdir pkg/starlark/fuzzdata -bin pkg/starlark/starlark-fuzz.zip

.PHONY: gen
