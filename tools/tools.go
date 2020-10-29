// +build tools

package tools // import "github.com/percona-platform/platform/tools"

import (
	// code generators plus their dependencies (to make them direct in go.mod)
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/mwitkow/go-proto-validators/protoc-gen-govalidators"
	_ "github.com/uber/prototool/cmd/prototool"
	_ "google.golang.org/grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"

	// other tools
	_ "github.com/dvyukov/go-fuzz/go-fuzz"
	_ "github.com/dvyukov/go-fuzz/go-fuzz-build"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/quasilyte/go-consistent"
	_ "github.com/reviewdog/reviewdog/cmd/reviewdog"
	_ "mvdan.cc/gofumpt/gofumports"
)
