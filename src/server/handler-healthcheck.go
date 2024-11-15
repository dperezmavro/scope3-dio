package server

import (
	"net/http"

	"github.com/scope3-dio/config"
)

func healthCheck(conf config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := healthCheckResponse{
			Environment: conf.Environment.Name,
			Service:     conf.Service.Name,
			Version:     conf.Service.Version,
		}
		writeResponse(w, r, body, http.StatusOK)
	}
}
