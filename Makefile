default: help

help:                 ## Display this help message.
	@echo "Please use \`make <target>\` where <target> is one of:"
	@grep '^[a-zA-Z]' $(MAKEFILE_LIST) | \
		awk -F ':.*?## ' 'NF==2 {printf "  %-26s%s\n", $$1, $$2}'


DOCKER_DEV_IMAGE  = percona-platform-prototool:dev
DOCKER_RUN_IMAGE ?= docker.pkg.github.com/percona-platform/platform/prototool:latest
DOCKER_RUN_CMD    = docker run --rm --mount='type=bind,src=$(PWD),dst=/work' $(DOCKER_RUN_IMAGE)


gen:                  ## Format, check, and generate using prototool Docker image.
	$(DOCKER_RUN_CMD) prototool break check api/telemetry -f api/telemetry/descriptor.bin

	rm -fr gen
	$(DOCKER_RUN_CMD) prototool all api
	$(DOCKER_RUN_CMD) gofmt -w -s .
	$(DOCKER_RUN_CMD) goimports -local github.com/percona-platform/platform -w .

test:
	go install ./...
	go test ./...

descriptors:          ## Update files used for breaking changes detection.
	$(DOCKER_RUN_CMD) prototool break descriptor-set api/telemetry -o api/telemetry/descriptor.bin

docker-build:         ## Build prototool Docker dev image.
	docker build --pull --squash --tag $(DOCKER_DEV_IMAGE) -f Dockerfile .

docker-push:          ## Tag and push prototool Docker image.
	docker tag $(DOCKER_DEV_IMAGE) $(DOCKER_RUN_IMAGE)
	docker push $(DOCKER_RUN_IMAGE)

run-dev:              ## Run bash in prototool Docker dev image.
	# the same as DOCKER_RUN_CMD but with `-it` and dev image
	docker run -it --rm --mount='type=bind,src=$(PWD),dst=/work' $(DOCKER_DEV_IMAGE) /bin/bash

saas:                 ## Extract public APIs and generated files into ../saas.
	rm -fr ../saas/api ../saas/gen
	mkdir ../saas/api ../saas/gen
	cp -R api/telemetry ../saas/api
	cp -R gen/telemetry ../saas/gen
	find ../saas -name '*.bin' -print -delete

ci:
	make docker-build
	env DOCKER_RUN_IMAGE=$(DOCKER_DEV_IMAGE) make gen

	go env
	sudo chown -R runner:docker gen
	go clean -testcache
	make test

	# Break CI if any files were changed during its run (code generation, etc), except go.sum.
	# `go mod tidy` could remove old checksums from that file, and that's okay on CI,
	# and actually expected for PRs made by @dependabot.
	# Checksums of actually used modules are checked by previous `go` subcommands.
	go mod tidy
	git checkout go.sum
	git diff --exit-code

.PHONY: help gen test docker-build docker-push run-dev saas ci
