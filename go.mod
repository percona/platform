module github.com/percona-platform/platform

go 1.13

// some dependecies should be synced with Dockerfile

require (
	github.com/golang/protobuf v1.3.4
	github.com/google/uuid v1.1.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/mwitkow/go-proto-validators v0.3.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.4.1
	github.com/stretchr/testify v1.5.1
	go.uber.org/zap v1.14.0
	golang.org/x/crypto v0.0.0-20200221231518-2aa609cf4a9d
	golang.org/x/sys v0.0.0-20200122134326-e047566fdf82
	google.golang.org/grpc v1.27.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)
