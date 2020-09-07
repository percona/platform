module github.com/percona-platform/platform/tools

go 1.14

// some dependecies should be synced with Dockerfile

require (
	github.com/dvyukov/go-fuzz v0.0.0-20200826052050-32ce4e791247
	github.com/elazarl/go-bindata-assetfs v1.0.0 // indirect
	github.com/golang/protobuf v1.3.5
	github.com/golangci/golangci-lint v1.29.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.7
	github.com/mwitkow/go-proto-validators v0.3.2
	github.com/quasilyte/go-consistent v0.0.0-20200404105227-766526bf1e96
	github.com/reviewdog/reviewdog v0.10.0
	github.com/stephens2424/writerset v1.0.2 // indirect
	github.com/uber/prototool v1.10.0
	google.golang.org/grpc v1.31.1
	mvdan.cc/gofumpt v0.0.0-20200802201014-ab5a8192947d
)
