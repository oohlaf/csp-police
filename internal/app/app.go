package app

import (
	"csp-police/pkg/api/csp"
	"github.com/rs/zerolog"
)

type appEnv struct {
	logger *zerolog.Logger
	client csp.CspServiceClient
}

func New(l *zerolog.Logger, c *csp.CspServiceClient) *appEnv {
	return &appEnv{
		logger: l,
		client: *c,
	}
}
