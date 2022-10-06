FROM golang:1.19-bullseye AS builder
WORKDIR /go/src/github.com/mkorman9/go-grpc-example
COPY . ./
RUN apt update && \
    apt install -y ca-certificates protobuf-compiler && \
    go install github.com/golang/protobuf/protoc-gen-go@latest && \
    make && \
    useradd -u 10001 app

FROM scratch
WORKDIR /
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/mkorman9/go-grpc-example/go-grpc-example .
USER app
ENV PROFILE=release
CMD ["./go-grpc-example"]
