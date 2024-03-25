module github.com/percona-platform/platform

go 1.21

// some dependecies should be synced with Dockerfile

require (
	github.com/Masterminds/squirrel v1.5.4
	github.com/alecthomas/kingpin/v2 v2.4.0
	github.com/aws/aws-sdk-go v1.51.7
	github.com/brianvoe/gofakeit/v6 v6.28.0
	github.com/golang/protobuf v1.5.4
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1
	github.com/lib/pq v1.9.0 // indirect
	github.com/mwitkow/go-proto-validators v0.3.2
	github.com/okta/okta-sdk-golang/v2 v2.20.0
	github.com/percona/promconfig v0.2.5
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.19.0
	github.com/stretchr/testify v1.9.0
	go.starlark.net v0.0.0-20230302034142-4b1e35fe2254
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.21.0
	golang.org/x/sys v0.18.0
	google.golang.org/genproto v0.0.0-20240123012728-ef4313101c80 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240125205218-1f4bbc51befe
	google.golang.org/grpc v1.62.1
	google.golang.org/protobuf v1.33.0
	gopkg.in/reform.v1 v1.5.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-jose/go-jose/v3 v3.0.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/patrickmn/go-cache v0.0.0-20180815053127-5633e0862627 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/xhit/go-str2duration/v2 v2.1.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240125205218-1f4bbc51befe // indirect
)
