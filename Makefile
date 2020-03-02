default: help

help:                 ## Display this help message.
	@echo "Please use \`make <target>\` where <target> is one of:"
	@grep '^[a-zA-Z]' $(MAKEFILE_LIST) | \
		awk -F ':.*?## ' 'NF==2 {printf "  %-26s%s\n", $$1, $$2}'


DOCKER_DEV_IMAGE  = percona-platform-prototool:dev
DOCKER_RUN_IMAGE ?= docker.pkg.github.com/percona-platform/platform/prototool:latest
DOCKER_RUN_CMD    = docker run --rm --mount='type=bind,src=$(PWD),dst=/work' $(DOCKER_RUN_IMAGE)


gen:                  ## Format, check, and generate using prototool Docker image.
	rm -fr gen

	# $(PROTOTOOL) break check api/events -f api/events/descriptor.bin
	$(DOCKER_RUN_CMD) prototool all api
	$(DOCKER_RUN_CMD) gofmt -w -s .
	$(DOCKER_RUN_CMD) goimports -local github.com/percona-platform/platform -w .

test:
	go install ./...
	go test ./...

# descriptors:          ## Update files used for breaking changes detection.
# 	$(PROTOTOOL) break descriptor-set api/events -o api/events/descriptor.bin
# 	$(PROTOTOOL) break descriptor-set api/callhome -o api/callhome/descriptor.bin

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

## Dummy

ci:
	make docker-build
	env DOCKER_RUN_IMAGE=$(DOCKER_DEV_IMAGE) make gen

	go env
	sudo chown -R runner:docker gen
	go clean -testcache
	make test
	go mod tidy
	git diff --exit-code

.PHONY: help gen test docker-build docker-push run-dev saas ci
