FROM debian:bullseye AS builder
RUN apt update && \
    apt install -y ca-certificates

FROM scratch
WORKDIR /
COPY go-grpc-example .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENV PROFILE=release
CMD ["./go-grpc-example"]
