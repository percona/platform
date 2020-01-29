FROM golang:1.13.7-buster

RUN apt-get update
RUN apt-get install -y aria2 unzip

# must match version in prototool.yml files
ENV PROTOBUF_VERSION=3.11.2

# must match versions in go.mod
ENV GRPC_GATEWAY_VERSION=1.12.1
ENV GO_PROTO_VALIDATORS_VERSION=0.3.0

RUN mkdir /tmp/protoc
RUN aria2c https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOBUF_VERSION}/protoc-${PROTOBUF_VERSION}-linux-x86_64.zip \
  --checksum=sha-512=50c639d8fed893acf28244f8119378b2d51918f7e24725d449c84d174ec5f6e71e939e58e42d60d86272e7bf638934855f5b03b03f7907b74b14225b924fd420 \
  --dir=/tmp/protoc --out=protoc.zip
RUN unzip /tmp/protoc/protoc.zip -d /tmp/protoc
RUN mv -v /tmp/protoc/include/* /usr/local/include
RUN mv -v /tmp/protoc/bin/* /usr/local/bin
RUN rm -frv /tmp/protoc

RUN mkdir /tmp/go
COPY go.mod go.sum tools.go /tmp/go/
RUN cd /tmp/go && env GO111MODULE=on go install -v -mod=readonly \
  github.com/gogo/protobuf/protoc-gen-gofast \
  github.com/gogo/protobuf/protoc-gen-gogo \
  github.com/gogo/protobuf/protoc-gen-gogofast \
  github.com/gogo/protobuf/protoc-gen-gogofaster \
  github.com/gogo/protobuf/protoc-gen-gogoslick \
  github.com/golang/protobuf/protoc-gen-go \
  github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
  github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
  github.com/mwitkow/go-proto-validators/protoc-gen-govalidators \
  github.com/uber/prototool/cmd/prototool \
  golang.org/x/tools/cmd/goimports
RUN mv -v /go/bin/* /usr/local/bin
RUN mv -v /go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v${GRPC_GATEWAY_VERSION}/third_party/googleapis/google/* /usr/local/include/google
RUN mkdir -p /usr/local/include/github.com/mwitkow/go-proto-validators
RUN mv -v /go/pkg/mod/github.com/mwitkow/go-proto-validators@v${GO_PROTO_VALIDATORS_VERSION}/*.proto /usr/local/include/github.com/mwitkow/go-proto-validators
RUN go clean -cache
RUN go clean -modcache
RUN rm -frv /go

ENV PROTOTOOL_PROTOC_BIN_PATH=/usr/local/bin/protoc
ENV PROTOTOOL_PROTOC_WKT_PATH=/usr/local/include

RUN protoc --version
RUN prototool version

WORKDIR /work
