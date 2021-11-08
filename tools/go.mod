module github.com/percona-platform/platform/tools

go 1.16

// some dependecies should be synced with Dockerfile

require (
	github.com/dvyukov/go-fuzz v0.0.0-20201127111758-49e582c6c23d
	github.com/elazarl/go-bindata-assetfs v1.0.1 // indirect
	github.com/fullstorydev/grpcurl v1.7.0 // indirect; https://github.com/uber/prototool/issues/559
	github.com/golang/protobuf v1.5.2
	github.com/golangci/golangci-lint v1.43.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/mwitkow/go-proto-validators v0.3.2
	github.com/quasilyte/go-consistent v0.0.0-20200404105227-766526bf1e96
	github.com/reviewdog/reviewdog v0.13.0
	github.com/stephens2424/writerset v1.0.2 // indirect
	github.com/uber/prototool v1.10.0
	golang.org/x/tools v0.1.7
	google.golang.org/grpc v1.41.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.27.1
	mvdan.cc/gofumpt v0.1.1
)
