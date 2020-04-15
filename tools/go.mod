module github.com/percona-platform/platform/tools

go 1.14

// some dependecies should be synced with Dockerfile

require (
	github.com/dvyukov/go-fuzz v0.0.0-20200318091601-be3528f3a813
	github.com/elazarl/go-bindata-assetfs v1.0.0 // indirect
	github.com/golang/protobuf v1.3.5
	github.com/golangci/golangci-lint v1.24.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.3
	github.com/mwitkow/go-proto-validators v0.3.0
	github.com/quasilyte/go-consistent v0.0.0-20200404105227-766526bf1e96
	github.com/reviewdog/reviewdog v0.9.17
	github.com/stephens2424/writerset v1.0.2 // indirect
	github.com/uber/prototool v1.9.0
	golang.org/x/tools v0.0.0-20200406213809-066fd1390ee0
	google.golang.org/grpc v1.28.1
)
