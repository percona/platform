module github.com/percona-platform/platform/tools

go 1.14

// some dependecies should be synced with Dockerfile

require (
	github.com/dvyukov/go-fuzz v0.0.0-20200826052050-32ce4e791247
	github.com/elazarl/go-bindata-assetfs v1.0.1 // indirect
	github.com/fullstorydev/grpcurl v1.7.0 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/golangci/golangci-lint v1.31.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.8
	github.com/jhump/protoreflect v1.7.0 // indirect
	github.com/mwitkow/go-proto-validators v0.3.2
	github.com/quasilyte/go-consistent v0.0.0-20200404105227-766526bf1e96
	github.com/reviewdog/reviewdog v0.10.2
	github.com/stephens2424/writerset v1.0.2 // indirect
	github.com/uber/prototool v1.10.0
	golang.org/x/net v0.0.0-20200904194848-62affa334b73 // indirect
	golang.org/x/sys v0.0.0-20200908134130-d2e65c121b96 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200904004341-0bd0a958aa1d // indirect
	google.golang.org/grpc v1.31.1
	google.golang.org/protobuf v1.25.0 // indirect
	mvdan.cc/gofumpt v0.0.0-20200802201014-ab5a8192947d
)
