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
	"strings"
)

const correctToken = "secret"

type GreeterService struct {
}

func (gs *GreeterService) SayHello(ctx context.Context, request *protocol.HelloRequest) (*protocol.HelloReply, error) {
	if !grpcutil.GetTokenVerificationResult(ctx).IsAuthorized {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	response := &protocol.HelloReply{
		Message: fmt.Sprintf("Hello %s", request.Name),
	}

	return response, nil
}

type GameService struct {
}

func (gs *GameService) Play(stream protocol.GameService_PlayServer) error {
	if !grpcutil.GetTokenVerificationResult(stream.Context()).IsAuthorized {
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	for {
		request, err := stream.Recv()
		if err != nil {
			return status.Errorf(codes.Canceled, "call cancelled")
		}

		if strings.EqualFold(request.Message, "exit") {
			_ = stream.Send(&protocol.ServerResponse{
				Message: "SERVERMSG: Closing",
			})

			break
		}

		_ = stream.Send(&protocol.ServerResponse{
			Message: fmt.Sprintf("SERVERMSG: Received %s", request.Message),
		})
	}

	return nil
}

func authFunction(token string) (*grpcutil.TokenVerificationResult, error) {
	return &grpcutil.TokenVerificationResult{
		IsAuthorized: token == correctToken,
	}, nil
}

func main() {
	configutil.LoadConfig()
	logutil.SetupLogger()

	server := grpcutil.NewServer(grpcutil.EnableAuthMiddleware(authFunction))
	protocol.RegisterGreeterServer(server.Server, &GreeterService{})
	protocol.RegisterGameServiceServer(server.Server, &GameService{})

	lifecycle.StartAndBlock(server)
}
