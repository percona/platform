// +build tools

package tools // import "github.com/percona-platform/platform/tools"

import (
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger"
	_ "github.com/mwitkow/go-proto-validators/protoc-gen-govalidators"
	_ "github.com/uber/prototool/cmd/prototool"
	_ "golang.org/x/tools/cmd/goimports"
)
