package main

import (
	"context"
	"fmt"
	"github.com/mkorman9/go-commons/configutil"
	"github.com/mkorman9/go-commons/coreutil"
	"github.com/mkorman9/go-commons/flyutil"
	"github.com/mkorman9/go-commons/grpcutil"
	"github.com/mkorman9/go-commons/logutil"
	"github.com/mkorman9/go-grpc-example/protocol"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

var AppVersion = "dev"

const correctToken = "secret"

type GreeterService struct {
}

func (gs *GreeterService) SayHello(ctx context.Context, request *protocol.HelloRequest) (*protocol.HelloReply, error) {
	if !grpcutil.GetTokenVerificationResult(ctx).IsAuthorized {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	response := &protocol.HelloReply{
		Message:     fmt.Sprintf("Hello %s from version %s", request.Name, AppVersion),
		InCloud:     flyutil.RunningInCloud(),
		Environment: coreutil.GetEnvironment(),
		Region:      flyutil.Region(),
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

	cancel := coreutil.Consume("server.events", func(msg string, _ coreutil.CallContext) {
		ds.Send(&protocol.ServerResponse{
			Message: msg,
		})
	})
	defer cancel()

	return ds.Start()
}

func authFunction(token string, _ *grpcutil.CallMetadata) (*grpcutil.TokenVerificationResult, error) {
	return &grpcutil.TokenVerificationResult{
		IsAuthorized: token == correctToken,
	}, nil
}

func main() {
	config := configutil.LoadConfig()
	logutil.SetupLogger(config)

	runScheduledEvents()

	server := grpcutil.NewServer(config, grpcutil.EnableAuthMiddlewareFunc(authFunction))
	protocol.RegisterGreeterServer(server.Server, &GreeterService{})
	protocol.RegisterGameServiceServer(server.Server, &GameService{})

	coreutil.StartAndBlock(server)
}

func runScheduledEvents() {
	go func() {
		for {
			coreutil.Publish(
				"server.events",
				"Scheduled event",
			)

			time.Sleep(5 * time.Second)
		}
	}()
}
