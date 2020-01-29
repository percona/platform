default: all

help:                 ## Display this help message.
	@echo "Please use \`make <target>\` where <target> is one of:"
	@grep '^[a-zA-Z]' $(MAKEFILE_LIST) | \
		awk -F ':.*?## ' 'NF==2 {printf "  %-26s%s\n", $$1, $$2}'

EXEC = docker run -it --rm --mount='type=bind,src=$(PWD),dst=/work' pp-builder:dev


all:                  ## Format, check, and generate.
	rm -fr gen

	# $(PROTOTOOL) break check api/events -f api/events/descriptor.bin
	$(EXEC) prototool all api
	$(EXEC) gofmt -w -s .
	$(EXEC) goimports -local github.com/Percona-Platform/platform -l -w .

# descriptors:          ## Update files used for breaking changes detection.
# 	$(PROTOTOOL) break descriptor-set api/events -o api/events/descriptor.bin
# 	$(PROTOTOOL) break descriptor-set api/callhome -o api/callhome/descriptor.bin

docker:
	docker build --pull --squash --tag pp-builder:dev -f Dockerfile .

exec:
	$(EXEC) /bin/bash
