package rest

import (
	"net/http"

	"metrics/internal/core/service"
	"metrics/internal/infra/api/rest/handlers"
)

type API struct {
	srv *http.ServeMux
}

func NewAPI(metricService *service.MetricService) *API {
	handler := handlers.NewHandler(metricService)

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handler.UpdateHandler)
	mux.HandleFunc("/value/", handler.GetHandler)

	return &API{
		srv: mux,
	}
}

func (app *API) Run() error {
	return http.ListenAndServe("localhost:8080", app.srv)
}
