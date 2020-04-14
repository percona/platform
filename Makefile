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
	go build -modfile=tools/go.mod -o bin/go-consistent github.com/quasilyte/go-consistent
	go build -modfile=tools/go.mod -o bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint
	go build -modfile=tools/go.mod -o bin/reviewdog github.com/reviewdog/reviewdog/cmd/reviewdog
	go build -modfile=tools/go.mod -o bin/go-fuzz github.com/dvyukov/go-fuzz/go-fuzz
	go build -modfile=tools/go.mod -o bin/go-fuzz-build github.com/dvyukov/go-fuzz/go-fuzz-build

gen:                                       ## Format, check, and generate using prototool Docker image
	$(DOCKER_RUN_CMD) prototool break check api/telemetry -f api/telemetry/descriptor.bin

	rm -rf gen
	$(DOCKER_RUN_CMD) prototool all api
	$(DOCKER_RUN_CMD) gofmt -w -s .
	$(DOCKER_RUN_CMD) goimports -local github.com/percona-platform/platform -w .

gen-dev: docker-build                      ## Same as `gen` but with DEV protocol Docker image
	env DOCKER_RUN_IMAGE=$(DOCKER_DEV_IMAGE) make gen
	sudo chown -R runner:docker gen

format:                                    ## Format source code
	gofmt -w -s .
	bin/goimports -local github.com/percona-platform/platform -l -w .

check:                                     ## Run checks/linters for the whole project
	bin/go-consistent -pedantic ./...
	bin/golangci-lint run

test:                                      ## Run tests
	go test ./...

descriptors:                               ## Update files used for breaking changes detection
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
	rm -rf ../saas/api ../saas/gen
	mkdir ../saas/api ../saas/gen
	cp -R api/telemetry ../saas/api
	cp -R gen/telemetry ../saas/gen
	find ../saas -name '*.bin' -print -delete

pkg/check/check-fuzz.zip:
	cd pkg/check && $(PWD)/bin/go-fuzz-build

fuzz-data: pkg/check/check-fuzz.zip        ## Fuzz data tests
	bin/go-fuzz -workdir pkg/check/fuzzdata -bin pkg/check/check-fuzz.zip -func FuzzData

fuzz-signature: pkg/check/check-fuzz.zip   ## Fuzz signature tests
	bin/go-fuzz -workdir pkg/check/fuzzdata -bin pkg/check/check-fuzz.zip -func FuzzSign

fuzz-pubkey: pkg/check/check-fuzz.zip      ## Fuzz public key tests
	bin/go-fuzz -workdir pkg/check/fuzzdata -bin pkg/check/check-fuzz.zip -func FuzzPublicKey

fuzz-clean:                                ## Clean fuzz data
	rm -rf pkg/check/fuzzdata
	rm -r pkg/check/check-fuzz.zip

.PHONY: gen
