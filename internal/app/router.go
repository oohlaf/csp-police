package app

import (
	"net/http"
)

func NewRouter(app *appEnv) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/csp-report", app.CreateCspReport)
	return mux
}
