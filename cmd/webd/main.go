package main

import (
	"fmt"
	"net/http"
	"os"

	"csp-police/internal/app"
	"csp-police/internal/config"
	"csp-police/pkg/api/csp"

	"google.golang.org/grpc"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	appConf := config.AppConfig()
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if appConf.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	grpcAddress := fmt.Sprintf("%s:%d", appConf.Grpc.Hostname, appConf.Grpc.Port)
	logger.Info().Str("address", grpcAddress).Msg("Dialing GRPC server")
	conn, err := grpc.Dial(grpcAddress, opts...)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to dial GRPC server")
	}
	defer conn.Close()
	client := csp.NewCspServiceClient(conn)

	application := app.New(&logger, &client)
	appRouter := app.NewRouter(application)
	webAddress := fmt.Sprintf("%s:%d", appConf.Web.Hostname, appConf.Web.Port)
	logger.Info().Str("address", webAddress).Msg("Starting webserver")

	if err := http.ListenAndServe(webAddress, appRouter); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start webserver")
	}
}
