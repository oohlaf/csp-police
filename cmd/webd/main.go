package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"csp-police/internal/config"
	"csp-police/internal/decode"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type BrowserCspReport struct {
	DocumentUri        string `json:"document-uri"`
	Referrer           string
	ViolatedDirective  string `json:"violated-directive"`
	EffectiveDirective string `json:"effective-directive"`
	OriginalPolicy     string `json:"original-policy"`
	Disposition        string
	BlockedUri         string `json:"blocked-uri"`
	StatusCode         int    `json:"status-code"`
	SourceFile         string `json:"source-file"`
	LineNumber         int    `json:"line-number"`
	ColumnNumber       int    `json:"column-number"`
	ScriptSample       string `json:"script-sample"`
}

type BrowserCspDocument struct {
	CspReport BrowserCspReport `json:"csp-report"`
}

func createCspReport(w http.ResponseWriter, r *http.Request) {
	var doc BrowserCspDocument

	if err := decode.DecodeJsonBody(w, r, &doc); err != nil {
		var mr *decode.MalformedRequest
		if errors.As(err, &mr) {
			log.Error().Err(err).Int("status", mr.Status()).Msg("")
			http.Error(w, mr.Error(), mr.Status())
		} else {
			log.Error().Err(err).Int("status", http.StatusInternalServerError).Msg(http.StatusText(http.StatusInternalServerError))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	log.Info().Msgf("CSP-Report %v", doc)
}

func main() {
	appConf := config.AppConfig()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if appConf.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/csp-report", createCspReport)
	address := fmt.Sprintf("%s:%d", appConf.Server.Hostname, appConf.Server.Port)
	log.Info().Str("address", address).Msg("Starting server")

	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatal().Err(err).Msg("Server startup failed")
	}
}
