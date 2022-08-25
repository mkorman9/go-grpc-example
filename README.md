# Prerequisites

- Protobuf compiler installed

macOS:
```
brew install protobuf
```

Linux (Debian-like):
```
sudo apt update && sudo apt install -y protobuf-compiler
```

- protoc-gen-go installed
```
go install github.com/golang/protobuf/protoc-gen-go@latest
```

(`$GOPATH/bin` must be in `$PATH`)

# Build

Build binary:
```
make generate
make 
```

Create Docker image:
```
docker build -t go-grpc-example .
```
