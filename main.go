package main

import (
	"context"
	"fmt"
	"github.com/mkorman9/go-commons/configutil"
	"github.com/mkorman9/go-commons/grpcutil"
	"github.com/mkorman9/go-commons/lifecycle"
	"github.com/mkorman9/go-commons/logutil"
	"github.com/mkorman9/go-grpc-example/protocol"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	ds := grpcutil.NewDuplexStream[protocol.PlayerRequest, protocol.ServerResponse](stream)

	ds.OnReceive(func(request *protocol.PlayerRequest) {
		if request.Message == "exit" {
			ds.Stop()
		} else {
			ds.Send(&protocol.ServerResponse{
				Message: request.Message,
			})
		}
	})

	ds.OnEnd(func(_ error) {
		log.Info().Msg("Connection closed")
	})

	return ds.Start()
}

func authFunction(token string, _ *grpcutil.CallMetadata) (*grpcutil.TokenVerificationResult, error) {
	return &grpcutil.TokenVerificationResult{
		IsAuthorized: token == correctToken,
	}, nil
}

func main() {
	configutil.LoadConfig()
	logutil.SetupLogger()

	server := grpcutil.NewServer(grpcutil.EnableAuthMiddlewareFunc(authFunction))
	protocol.RegisterGreeterServer(server.Server, &GreeterService{})
	protocol.RegisterGameServiceServer(server.Server, &GameService{})

	lifecycle.StartAndBlock(server)
}
