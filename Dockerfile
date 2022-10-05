FROM golang:1.19-bullseye AS builder
WORKDIR /go/src/github.com/mkorman9/go-grpc-example
RUN apt update && \
    apt install -y ca-certificates protobuf-compiler && \
    go install github.com/golang/protobuf/protoc-gen-go@latest
COPY . ./
RUN make

FROM scratch
WORKDIR /
COPY go-grpc-example .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/mkorman9/go-grpc-example/go-grpc-example .
ENV PROFILE=release
CMD ["./go-grpc-example"]
