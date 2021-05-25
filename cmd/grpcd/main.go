package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"csp-police/internal/config"
	"csp-police/pkg/api/csp"

	"google.golang.org/grpc"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type cspServiceServer struct {
	csp.UnimplementedCspServiceServer
}

func (s *cspServiceServer) Send(ctx context.Context, request *csp.SendRequest) (*csp.SendResponse, error) {
	log.Info().Str("agent", request.Client.UserAgent).Msg("Received data")
	return &csp.SendResponse{}, nil
}

func NewServer() *cspServiceServer {
	return &cspServiceServer{}
}

func main() {
	appConf := config.AppConfig()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if appConf.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	address := fmt.Sprintf("%s:%d", appConf.Grpc.Hostname, appConf.Grpc.Port)
	log.Info().Str("address", address).Msg("Starting GRPC server")

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal().Err(err).Msg("GRPC server failed to listen")
	}

	s := grpc.NewServer()
	csp.RegisterCspServiceServer(s, NewServer())
	if err := s.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("Failed to start GRPC server")
	}
}
