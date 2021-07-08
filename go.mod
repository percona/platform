module github.com/percona-platform/platform

go 1.16

// some dependecies should be synced with Dockerfile

require (
	github.com/golang/protobuf v1.5.1
	github.com/google/uuid v1.2.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/mwitkow/go-proto-validators v0.3.2
	github.com/percona/promconfig v0.2.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/stretchr/testify v1.7.0
	go.starlark.net v0.0.0-20210312235212-74c10e2c17dc
	go.uber.org/zap v1.18.1
	golang.org/x/sys v0.0.0-20210603081109-ebe580a85c40
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/reform.v1 v1.5.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)
