FROM golang:1.16.6

RUN apt-get update
RUN apt-get install -y aria2 unzip

# must match version in prototool.yml file
ENV PROTOBUF_VERSION=3.15.6
ENV PROTOBUF_CHECKSUM=1c7c98819985c3d5284bb9baf423cf701a678372a46f7a0fd7c48dee398847032c4727dd32713ad99517e63fe1fb59976b4f46ebc3b8bef0bc14d9a9180f4111

# must match versions in tools/go.mod
ENV GO_PROTO_VALIDATORS_VERSION=0.3.2

ENV GRPC_WEB_VERSION=1.2.1
ENV GRPC_WEB_CHECKSUM=6ce1625db7902d38d38d83690ec578c182e9cf2abaeb58d3fba1dae0c299c597

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
  github.com/mwitkow/go-proto-validators/protoc-gen-govalidators \
  github.com/uber/prototool/cmd/prototool \
  mvdan.cc/gofumpt/gofumports
RUN mv -v /go/bin/* /usr/local/bin
RUN mkdir -p /usr/local/include/github.com/mwitkow/go-proto-validators
RUN mv -v /go/pkg/mod/github.com/mwitkow/go-proto-validators@v${GO_PROTO_VALIDATORS_VERSION}/*.proto /usr/local/include/github.com/mwitkow/go-proto-validators
RUN go clean -cache
RUN go clean -modcache
RUN rm -frv /go

RUN aria2c https://github.com/grpc/grpc-web/releases/download/${GRPC_WEB_VERSION}/protoc-gen-grpc-web-${GRPC_WEB_VERSION}-linux-x86_64 \
  --checksum=sha-256=${GRPC_WEB_CHECKSUM} \
  --dir /usr/local/bin --out=protoc-gen-grpc-web && \
  chmod +x /usr/local/bin/protoc-gen-grpc-web

RUN curl -sL https://deb.nodesource.com/setup_12.x | bash -
RUN apt-get install -y nodejs
RUN npm install --global --unsafe-perm --production grpc-tools

ENV PROTOTOOL_PROTOC_BIN_PATH=/usr/local/bin/protoc
ENV PROTOTOOL_PROTOC_WKT_PATH=/usr/local/include

RUN protoc --version
RUN prototool version

WORKDIR /work
