module github.com/percona-platform/platform/tools

go 1.14

// some dependecies should be synced with Dockerfile

require (
	github.com/dvyukov/go-fuzz v0.0.0-20200318091601-be3528f3a813
	github.com/elazarl/go-bindata-assetfs v1.0.0 // indirect
	github.com/golang/protobuf v1.3.5
	github.com/golangci/golangci-lint v1.28.1
	github.com/grpc-ecosystem/grpc-gateway v1.14.6
	github.com/mwitkow/go-proto-validators v0.3.0
	github.com/quasilyte/go-consistent v0.0.0-20200404105227-766526bf1e96
	github.com/reviewdog/reviewdog v0.10.0
	github.com/stephens2424/writerset v1.0.2 // indirect
	github.com/uber/prototool v1.10.0
	golang.org/x/tools v0.0.0-20200702044944-0cc1aa72b347
	google.golang.org/grpc v1.30.0
)
