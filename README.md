# platform

[![CI GitHub Action status](https://github.com/percona-platform/platform/workflows/CI/badge.svg?branch=main)](https://github.com/percona-platform/platform/actions?query=workflow%3ACI+branch%3Amain)

Percona Platform private API definitions and common code.

## Contributing

- All the API definitions are contained inside `/api` directory.
- `/pkg` contains common code used by other platform projects.
- Run `make init` to install development tools, and `make gen` to format, check, and generate code from API definitions.
- Run `make` to see all targets.
- Please run `make test` and `make check` before submitting a PR.

### Working with Buf

We use [buf](https://buf.build/) to generate code from our protobuf files. For more information take a look at their [docs](https://docs.buf.build/introduction)

The `/api` directory has a `buf.yaml` file to make our protofiles discoverable to buf. This file also contains configuration for break checks, linting and 3rd party protobuf dependencies. These dependencies are downloaded from Buf Schema Registry when you run buf for the first time. If any dependencies are declared in the `buf.yaml` file, it is accompanied by a `buf.lock` file that contains manifests for the listed dependencies. `buf.lock` can only be created/updated manually using `buf mod update` command. Downloading dependencies from Buf Schema Registry does not require any kind of extra configuration or authentication, you do not need to have a schema registry account for it. These dependencies are downloaded and stored in a [module cache](https://docs.buf.build/bsr/overview#module-cache) locally to avoid multiple calls to BSR.

By default buf downloads the latest version of the dependency, if you want a specific version you can pin a dependency using a commit hash or tag. [More info here](https://docs.buf.build/tour/add-a-dependency#pin-your-dependencies)

The `buf.gen.yaml` file defines code generation options for all the protoc plugins used by the project. The `buf.work.yaml` file defines the buf workspace. Buf generates a single descriptor binary (named `platform.bin`) as opposed to individual descriptors for each API.
