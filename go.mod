module github.com/percona-platform/platform

go 1.14

// some dependecies should be synced with Dockerfile

require (
	github.com/AlekSi/pointer v1.1.0 // indirect
	github.com/denisenkom/go-mssqldb v0.0.0-20200206145737-bbfc9a55622e // indirect
	github.com/dvyukov/go-fuzz v0.0.0-20200318091601-be3528f3a813 // indirect
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/golang/protobuf v1.3.5
	github.com/google/uuid v1.1.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/lib/pq v1.4.0 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/mwitkow/go-proto-validators v0.3.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.5.1
	github.com/stretchr/testify v1.5.1
	go.starlark.net v0.0.0-20200330013621-be5394c419b6
	go.uber.org/zap v1.15.0
	golang.org/x/sys v0.0.0-20200427175716-29b57079015a
	google.golang.org/grpc v1.29.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/reform.v1 v1.3.3
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
	syreclabs.com/go/faker v1.2.2 // indirect
)
