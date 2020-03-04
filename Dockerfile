# must match version in prototool.yml files
ARG PROTOBUF_VERSION=3.11.2
ARG PROTOBUF_CHECKSUM=50c639d8fed893acf28244f8119378b2d51918f7e24725d449c84d174ec5f6e71e939e58e42d60d86272e7bf638934855f5b03b03f7907b74b14225b924fd420

# must match versions in tools/go.mod
ARG GRPC_GATEWAY_VERSION=1.13.0
ARG GO_PROTO_VALIDATORS_VERSION=0.3.0

FROM golang:1.13 as build

WORKDIR /tmp/build

ARG PROTOBUF_VERSION
ARG PROTOBUF_CHECKSUM

RUN apt-get update && apt-get install -y aria2 unzip
RUN mkdir /tmp/protoc && \
    aria2c https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOBUF_VERSION}/protoc-${PROTOBUF_VERSION}-linux-x86_64.zip \
        --checksum=sha-512=${PROTOBUF_CHECKSUM} --out=protoc.zip && \
    unzip protoc.zip -d /tmp/protoc

COPY tools/go.mod tools/go.sum tools/tools.go ./
RUN go mod download
RUN go install -v -mod=readonly \
        github.com/golang/protobuf/protoc-gen-go \
        github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
        github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
        github.com/mwitkow/go-proto-validators/protoc-gen-govalidators \
        github.com/uber/prototool/cmd/prototool \
        golang.org/x/tools/cmd/goimports


FROM golang:1.13 as target

WORKDIR /work

ARG GRPC_GATEWAY_VERSION
ARG GO_PROTO_VALIDATORS_VERSION

ENV PROTOTOOL_PROTOC_BIN_PATH=/usr/local/bin/protoc
ENV PROTOTOOL_PROTOC_WKT_PATH=/usr/local/include

COPY --from=build /go/bin /tmp/protoc/bin /usr/local/bin/
COPY --from=build /tmp/protoc/include /usr/local/include
COPY --from=build /go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v${GRPC_GATEWAY_VERSION}/third_party/googleapis/google /usr/local/include/google
COPY --from=build /go/pkg/mod/github.com/mwitkow/go-proto-validators@v${GO_PROTO_VALIDATORS_VERSION}/*.proto /usr/local/include/github.com/mwitkow/go-proto-validators/
