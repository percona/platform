module github.com/percona-platform/platform/tools

go 1.15

// some dependecies should be synced with Dockerfile

require (
	github.com/dvyukov/go-fuzz v0.0.0-20201003075337-90825f39c90b
	github.com/elazarl/go-bindata-assetfs v1.0.1 // indirect
	github.com/golang/protobuf v1.4.3
	github.com/golangci/golangci-lint v1.33.0
	github.com/mwitkow/go-proto-validators v0.3.2
	github.com/quasilyte/go-consistent v0.0.0-20200404105227-766526bf1e96
	github.com/reviewdog/reviewdog v0.11.0
	github.com/stephens2424/writerset v1.0.2 // indirect
	github.com/uber/prototool v1.10.0
	golang.org/x/tools v0.0.0-20201105220310-78b158585360
	google.golang.org/grpc v1.34.0
	google.golang.org/protobuf v1.25.0
	mvdan.cc/gofumpt v0.0.0-20201107090320-a024667a00f1
)
