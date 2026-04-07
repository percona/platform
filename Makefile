default: help

help:                                      ## Display this help message
	@echo "Please use \`make <target>\` where <target> is one of:"
	@grep '^[a-zA-Z]' $(MAKEFILE_LIST) | \
		awk -F ':.*?## ' 'NF==2 {printf "  %-26s%s\n", $$1, $$2}'

init:                                      ## Install development tools
	rm -rf bin
	cd tools && go generate -x -tags=tools

ci-init:                                   ## Initialize CI environment

gen:                                       ## Format, check, and generate code using buf; TODO Add lint and break commands
	rm -rf gen
	bin/buf generate -v api
	make format
	bin/buf breaking --against platform.bin api

gen-code:                                  ## Generate code
	go generate ./...
	go install ./...

swagger-ui:                                ## Serve API documentation with SwaggerUI
	docker run -p 8081:8080 -e URLS='[ \
		{name:"telemetryd", url:"/gen/telemetry/reporter/reporter_api.swagger.json"}, \
		{name:"generic-telemetryd", url:"/gen/telemetry/generic/reporter_api.swagger.json"}, \
]' -v ./gen:/usr/share/nginx/html/gen swaggerapi/swagger-ui

format:                                    ## Format source code
	bin/gofumpt -l -w .
	bin/goimports -local github.com/percona/platform -l -w .
	bin/buf format api -w

check:                                     ## Run checks/linters for the whole project
	bin/go-consistent -exclude=tools -pedantic ./...
	LOG_LEVEL=error bin/golangci-lint run

test:                                      ## Run tests
	go test -race ./...

test-cover:                                ## Run tests and collect per-package coverage information
	go test -race -timeout=10m -count=1 -coverprofile=cover.out -covermode=atomic ./...

test-crosscover:                           ## Run tests and collect cross-package coverage information
	go test -race -timeout=10m -count=1 -coverprofile=crosscover.out -covermode=atomic -p=1 -coverpkg=./... ./...

descriptors:                               ## Update files used for breaking changes detection
	bin/buf build -o platform.bin --as-file-descriptor-set api

saas:                                      ## Extract public APIs and generated files into ../saas
	go run post-processing.go -project saas

.PHONY: $(MAKECMDGOALS)
