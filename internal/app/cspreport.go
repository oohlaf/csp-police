package app

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"csp-police/internal/decode"
	"csp-police/pkg/api/csp"

	"github.com/rs/zerolog"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type BrowserCspReport struct {
	DocumentUri        *string `json:"document-uri,omitempty"`
	Referrer           *string `json:"referrer,omitempty"`
	ViolatedDirective  *string `json:"violated-directive,omitempty"`
	EffectiveDirective *string `json:"effective-directive,omitempty"`
	OriginalPolicy     *string `json:"original-policy,omitempty"`
	Disposition        *string `json:"disposition,omitempty"`
	BlockedUri         *string `json:"blocked-uri,omitempty"`
	StatusCode         *int    `json:"status-code,omitempty"`
	SourceFile         *string `json:"source-file,omitempty"`
	LineNumber         *int    `json:"line-number,omitempty"`
	ColumnNumber       *int    `json:"column-number,omitempty"`
	ScriptSample       *string `json:"script-sample,omitempty"`
}

type BrowserCspDocument struct {
	CspReport BrowserCspReport `json:"csp-report"`
}

func (r BrowserCspReport) MarshalZerologObject(e *zerolog.Event) {
	e.Str("document-uri", NullString(r.DocumentUri)).
		Str("referrer", NullString(r.Referrer)).
		Str("violated-directive", NullString(r.ViolatedDirective)).
		Str("effective-directive", NullString(r.EffectiveDirective)).
		Str("original-policy", NullString(r.OriginalPolicy)).
		Str("disposition", NullString(r.Disposition)).
		Str("blocked-uri", NullString(r.BlockedUri)).
		Str("status-code", NullIntToString(r.StatusCode)).
		Str("source-file", NullString(r.SourceFile)).
		Str("line-number", NullIntToString(r.LineNumber)).
		Str("column-number", NullIntToString(r.ColumnNumber)).
		Str("script-sample", NullString(r.ScriptSample))
}

func NullString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func NullIntToString(value *int) string {
	if value == nil {
		return ""
	}
	return strconv.Itoa(*value)
}

func StringToBool(value string) bool {
	if strings.Compare(value, "") == 0 {
		return false
	}
	if strings.Compare(value, "0") == 0 {
		return false
	}
	return true
}

func (app *appEnv) CreateCspReport(w http.ResponseWriter, r *http.Request) {
	var doc BrowserCspDocument

	if err := decode.DecodeJsonBody(w, r, &doc); err != nil {
		var mr *decode.MalformedRequest
		if errors.As(err, &mr) {
			app.logger.Error().Err(err).Int("status", mr.Status()).Msg("")
			http.Error(w, mr.Error(), mr.Status())
		} else {
			app.logger.Error().Err(err).Int("status", http.StatusInternalServerError).Msg(http.StatusText(http.StatusInternalServerError))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	app.logger.Debug().Object("csp-report", doc.CspReport).Msg("Received CSP report")
	param := csp.Param{
		Application: r.FormValue("a"),
		Version:     r.FormValue("v"),
		Checksum:    r.FormValue("c"),
		Force:       StringToBool(r.FormValue("f")),
	}
	client := csp.Client{
		RemoteAddress: r.RemoteAddr,
		UserAgent:     r.UserAgent(),
	}
	report := csp.CspReport{
		DocumentUri:        NullString(doc.CspReport.DocumentUri),
		ReferrerUri:        NullString(doc.CspReport.Referrer),
		ViolatedDirective:  NullString(doc.CspReport.ViolatedDirective),
		EffectiveDirective: NullString(doc.CspReport.EffectiveDirective),
		OriginalPolicy:     NullString(doc.CspReport.OriginalPolicy),
		Disposition:        NullString(doc.CspReport.Disposition),
		BlockedUri:         NullString(doc.CspReport.BlockedUri),
		StatusCode:         NullIntToString(doc.CspReport.StatusCode),
		SourceUri:          NullString(doc.CspReport.SourceFile),
		LineNumber:         NullIntToString(doc.CspReport.LineNumber),
		ColumnNumber:       NullIntToString(doc.CspReport.ColumnNumber),
		ScriptSample:       NullString(doc.CspReport.ScriptSample),
	}
	sendRequest := csp.SendRequest{
		Timestamp: timestamppb.Now(),
		Param:     &param,
		Client:    &client,
		Report:    &report,
	}
	app.client.Send(context.Background(), &sendRequest)
}
