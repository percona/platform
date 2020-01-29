default: help

help:                 ## Display this help message.
	@echo "Please use \`make <target>\` where <target> is one of:"
	@grep '^[a-zA-Z]' $(MAKEFILE_LIST) | \
		awk -F ':.*?## ' 'NF==2 {printf "  %-26s%s\n", $$1, $$2}'


EXEC = docker run -it --rm --mount='type=bind,src=$(PWD),dst=/work' docker.pkg.github.com/percona-platform/platform/prototool:latest


all:                  ## Format, check, and generate using prototool Docker image.
	rm -fr gen

	# $(PROTOTOOL) break check api/events -f api/events/descriptor.bin
	$(EXEC) prototool all api
	$(EXEC) gofmt -w -s .
	$(EXEC) goimports -local github.com/Percona-Platform/platform -l -w .

# descriptors:          ## Update files used for breaking changes detection.
# 	$(PROTOTOOL) break descriptor-set api/events -o api/events/descriptor.bin
# 	$(PROTOTOOL) break descriptor-set api/callhome -o api/callhome/descriptor.bin

docker:               ## Build prototool Docker image.
	docker build --pull --squash --tag percona-platform-prototool:dev -f Dockerfile .
	docker tag percona-platform-prototool:dev docker.pkg.github.com/percona-platform/platform/prototool:latest
	docker push docker.pkg.github.com/percona-platform/platform/prototool:latest

exec:                 ## Run bash in prototool Docker image.
	$(EXEC) /bin/bash
