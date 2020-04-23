FROM golang:1.14

RUN apt-get update
RUN apt-get install -y aria2 unzip

# must match version in prototool.yml file
ENV PROTOBUF_VERSION=3.11.4
ENV PROTOBUF_CHECKSUM=82777f04d9600ec69c53044a06fec4d3e108c9c3797d643f3472eb558088963e02a153077e2f832db54d17921204d327ad6ba9f37db7d00bd46f4887229dc837

# must match versions in tools/go.mod
ENV GRPC_GATEWAY_VERSION=1.14.4
ENV GO_PROTO_VALIDATORS_VERSION=0.3.0

RUN mkdir /tmp/protoc
RUN echo https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOBUF_VERSION}/protoc-${PROTOBUF_VERSION}-linux-x86_64.zip
RUN aria2c https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOBUF_VERSION}/protoc-${PROTOBUF_VERSION}-linux-x86_64.zip \
  --checksum=sha-512=${PROTOBUF_CHECKSUM} \
  --dir=/tmp/protoc --out=protoc.zip
RUN unzip /tmp/protoc/protoc.zip -d /tmp/protoc
RUN mv -v /tmp/protoc/include/* /usr/local/include
RUN mv -v /tmp/protoc/bin/* /usr/local/bin
RUN rm -frv /tmp/protoc

RUN mkdir /tmp/go
COPY tools/go.mod tools/go.sum tools/tools.go /tmp/go/
RUN cd /tmp/go && go install -v -mod=readonly \
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
