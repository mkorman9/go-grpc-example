package main

import (
	"context"
	"fmt"
	"github.com/mkorman9/go-commons/configutil"
	"github.com/mkorman9/go-commons/grpcutil"
	"github.com/mkorman9/go-commons/lifecycle"
	"github.com/mkorman9/go-commons/logutil"
	"github.com/mkorman9/go-grpc-example/protocol"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GreeterService struct {
}

func (gs *GreeterService) SayHello(ctx context.Context, request *protocol.HelloRequest) (*protocol.HelloReply, error) {
	token := grpcutil.GetToken(ctx)
	if token != "secret" {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	response := &protocol.HelloReply{
		Message: fmt.Sprintf("Hello %s", request.Name),
	}

	return response, nil
}

func main() {
	configutil.LoadConfig()
	logutil.SetupLogger()

	server := grpcutil.NewServer()
	protocol.RegisterGreeterServer(server.Server, &GreeterService{})

	lifecycle.StartAndBlock(server)
}
