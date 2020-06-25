module github.com/percona-platform/platform

go 1.14

// some dependecies should be synced with Dockerfile

require (
	github.com/AlekSi/pointer v1.1.0 // indirect
	github.com/brianvoe/gofakeit v3.18.0+incompatible // indirect
	github.com/denisenkom/go-mssqldb v0.0.0-20200620013148-b91950f658ec // indirect
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/lib/pq v1.7.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.0 // indirect
	github.com/mwitkow/go-proto-validators v0.3.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.5.1
	github.com/stretchr/testify v1.6.1
	go.starlark.net v0.0.0-20200619143648-50ca820fafb9
	go.uber.org/zap v1.15.0
	golang.org/x/sys v0.0.0-20200622214017-ed371f2e16b4
	google.golang.org/grpc v1.30.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/reform.v1 v1.3.4
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)
