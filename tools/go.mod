module github.com/percona-platform/platform/tools

go 1.14

// some dependecies should be synced with Dockerfile

require (
	github.com/golang/protobuf v1.3.5
	github.com/golangci/golangci-lint v1.24.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.13.0
	github.com/mwitkow/go-proto-validators v0.3.0
	github.com/quasilyte/go-consistent v0.0.0-20190521200055-c6f3937de18c // indirect
	github.com/reviewdog/reviewdog v0.9.17 // indirect
	github.com/uber/prototool v1.9.0
	golang.org/x/tools v0.0.0-20200221224223-e1da425f72fd
)
