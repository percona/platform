DOCKER_DEV_IMAGE  = percona-platform-prototool:dev
DOCKER_RUN_IMAGE ?= docker.pkg.github.com/percona-platform/platform/prototool:latest
DOCKER_RUN_CMD    = docker run --rm --mount='type=bind,src=$(PWD),dst=/work' $(DOCKER_RUN_IMAGE)

default: help

help:                                      ## Display this help message
	@echo "Please use \`make <target>\` where <target> is one of:"
	@grep '^[a-zA-Z]' $(MAKEFILE_LIST) | \
		awk -F ':.*?## ' 'NF==2 {printf "  %-26s%s\n", $$1, $$2}'

init:                                      ## Install development tools
	go build -modfile=tools/go.mod -o bin/stringer golang.org/x/tools/cmd/stringer
	go build -modfile=tools/go.mod -o bin/go-consistent github.com/quasilyte/go-consistent
	go build -modfile=tools/go.mod -o bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint
	go build -modfile=tools/go.mod -o bin/reviewdog github.com/reviewdog/reviewdog/cmd/reviewdog
	go build -modfile=tools/go.mod -o bin/go-fuzz github.com/dvyukov/go-fuzz/go-fuzz
	go build -modfile=tools/go.mod -o bin/go-fuzz-build github.com/dvyukov/go-fuzz/go-fuzz-build
	go build -modfile=tools/go.mod -o bin/gofumports mvdan.cc/gofumpt/gofumports

ci-init:                                   ## Initialize CI environment
	# nothing there yet

gen:                                       ## Format, check, and generate code using prototool Docker image
	$(DOCKER_RUN_CMD) prototool break check api/auth -f api/auth/descriptor.bin
	$(DOCKER_RUN_CMD) prototool break check api/check/retrieval -f api/check/retrieval/descriptor.bin
	$(DOCKER_RUN_CMD) prototool break check api/telemetry -f api/telemetry/descriptor.bin

	rm -rf gen
	$(DOCKER_RUN_CMD) prototool all api
	$(DOCKER_RUN_CMD) gofumports -local github.com/percona-platform/platform -w .

	$(DOCKER_RUN_CMD) go run post-processing.go -patch-ui

gen-dev: docker-build                      ## Same as `gen` but with DEV prototool Docker image
	env DOCKER_RUN_IMAGE=$(DOCKER_DEV_IMAGE) make gen

gen-code:                                  ## Generate code
	go generate ./...
	go install ./...

format:                                    ## Format source code
	bin/gofumports -local github.com/percona-platform/platform -l -w .

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
	$(DOCKER_RUN_CMD) prototool break descriptor-set api/auth -o api/auth/descriptor.bin
	$(DOCKER_RUN_CMD) prototool break descriptor-set api/check/retrieval -o api/check/retrieval/descriptor.bin
	$(DOCKER_RUN_CMD) prototool break descriptor-set api/telemetry -o api/telemetry/descriptor.bin

docker-build:                              ## Build prototool Docker dev image
	docker build --pull --squash --tag $(DOCKER_DEV_IMAGE) -f Dockerfile .

docker-push:                               ## Tag and push prototool Docker image
	docker tag $(DOCKER_DEV_IMAGE) $(DOCKER_RUN_IMAGE)
	docker push $(DOCKER_RUN_IMAGE)

run-dev:                                   ## Run bash in prototool Docker dev image
	# the same as DOCKER_RUN_CMD but with `-it` and dev image
	docker run -it --rm --mount='type=bind,src=$(PWD),dst=/work' $(DOCKER_DEV_IMAGE) /bin/bash

saas:                                      ## Extract public APIs and generated files into ../saas
	go run post-processing.go -project saas

saas-ui:                                   ## Extract generated JS/TS files into ../saas-ui
	go run post-processing.go -project saas-ui

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
