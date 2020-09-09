module github.com/percona-platform/platform

go 1.14

// some dependecies should be synced with Dockerfile

require (
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.1
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/mwitkow/go-proto-validators v0.3.2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.7.1
	github.com/stretchr/testify v1.6.1
	go.starlark.net v0.0.0-20200901195727-6e684ef5eeee
	go.uber.org/zap v1.16.0
	golang.org/x/sys v0.0.0-20200908134130-d2e65c121b96
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/reform.v1 v1.4.1
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)
