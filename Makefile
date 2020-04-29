DOCKER_DEV_IMAGE  = percona-platform-prototool:dev
DOCKER_RUN_IMAGE ?= docker.pkg.github.com/percona-platform/platform/prototool:latest
DOCKER_RUN_CMD    = docker run --rm --mount='type=bind,src=$(PWD),dst=/work' $(DOCKER_RUN_IMAGE)

default: help

help:                                      ## Display this help message
	@echo "Please use \`make <target>\` where <target> is one of:"
	@grep '^[a-zA-Z]' $(MAKEFILE_LIST) | \
		awk -F ':.*?## ' 'NF==2 {printf "  %-26s%s\n", $$1, $$2}'

init:                                      ## Install development tools
	go build -modfile=tools/go.mod -o bin/goimports golang.org/x/tools/cmd/goimports
	go build -modfile=tools/go.mod -o bin/stringer golang.org/x/tools/cmd/stringer
	go build -modfile=tools/go.mod -o bin/go-consistent github.com/quasilyte/go-consistent
	go build -modfile=tools/go.mod -o bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint
	go build -modfile=tools/go.mod -o bin/reviewdog github.com/reviewdog/reviewdog/cmd/reviewdog
	go build -modfile=tools/go.mod -o bin/go-fuzz github.com/dvyukov/go-fuzz/go-fuzz
	go build -modfile=tools/go.mod -o bin/go-fuzz-build github.com/dvyukov/go-fuzz/go-fuzz-build

gen:                                       ## Format, check, and generate using prototool Docker image
	$(DOCKER_RUN_CMD) prototool break check api/telemetry -f api/telemetry/descriptor.bin
	$(DOCKER_RUN_CMD) prototool break check api/check/retrieval -f api/check/retrieval/descriptor.bin

	rm -rf gen
	$(DOCKER_RUN_CMD) prototool all api
	$(DOCKER_RUN_CMD) gofmt -w -s .
	$(DOCKER_RUN_CMD) goimports -local github.com/percona-platform/platform -w .

gen-dev: docker-build                      ## Same as `gen` but with DEV protocol Docker image
	env DOCKER_RUN_IMAGE=$(DOCKER_DEV_IMAGE) make gen
	sudo chown -R runner:docker gen

gen-code:                                 ## Generate code
	go generate ./...

format:                                    ## Format source code
	gofmt -w -s .
	bin/goimports -local github.com/percona-platform/platform -l -w .

check:                                     ## Run checks/linters for the whole project
	bin/go-consistent -pedantic ./...
	bin/golangci-lint run

test:                                      ## Run tests
	go test -race ./...

descriptors:                               ## Update files used for breaking changes detection
	$(DOCKER_RUN_CMD) prototool break descriptor-set api/telemetry -o api/telemetry/descriptor.bin
	$(DOCKER_RUN_CMD) prototool break descriptor-set api/check/retrieval -o api/check/retrieval/descriptor.bin

docker-build:                              ## Build prototool Docker dev image
	docker build --pull --squash --tag $(DOCKER_DEV_IMAGE) -f Dockerfile .

docker-push:                               ## Tag and push prototool Docker image
	docker tag $(DOCKER_DEV_IMAGE) $(DOCKER_RUN_IMAGE)
	docker push $(DOCKER_RUN_IMAGE)

run-dev:                                   ## Run bash in prototool Docker dev image
	# the same as DOCKER_RUN_CMD but with `-it` and dev image
	docker run -it --rm --mount='type=bind,src=$(PWD),dst=/work' $(DOCKER_DEV_IMAGE) /bin/bash

saas:                                      ## Extract public APIs and generated files into ../saas
	rm -rf ../saas/api ../saas/gen ../saas/pkg
	mkdir ../saas/api ../saas/gen ../saas/pkg

	cp -R api/check/retrieval ../saas/api
	cp -R api/telemetry ../saas/api

	cp -R gen/check/retrieval ../saas/gen
	cp -R gen/telemetry ../saas/gen

	cp -R pkg/check ../saas/pkg
	cp -R pkg/starlark ../saas/pkg

	gofmt -s -w -r '"github.com/percona-platform/platform/pkg/check" -> "github.com/percona-platform/saas/pkg/check"' ../saas/pkg

	find ../saas -name '*.bin' -print -delete
	find ../saas -name '*_test.go' -print -delete

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
