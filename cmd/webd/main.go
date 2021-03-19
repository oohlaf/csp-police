package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/oohlaf/csp-testing/internal/decode"
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
	err := decode.DecodeJsonBody(w, r, &doc)
	if err != nil {
		var mr *decode.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Error(), mr.Status())
		} else {
			// log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintf(w, "CSP-report: %+v", doc)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/csp-report", createCspReport)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Print(err)
	}
	// log.Fatal(err)
}
