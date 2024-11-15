package server

import (
	"net/http"

	"github.com/scope3-dio/config"
)

func measure(conf config.Config) http.HandlerFunc {
	return ResponseWritter("UNIMPLEMENTED", http.StatusInternalServerError)
}

func healthCheck(conf config.Config) http.HandlerFunc {
	body := healthCheckResponse{
		Environment: conf.Environment.Name,
		Service:     conf.Service.Name,
		Version:     conf.Service.Version,
	}
	return ResponseWritter(body, http.StatusOK)
}
