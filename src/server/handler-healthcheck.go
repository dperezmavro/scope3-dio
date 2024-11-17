package server

import (
	"net/http"

	"github.com/dperezmavro/scope3-dio/src/config"
	"github.com/dperezmavro/scope3-dio/src/utils"
)

// healthCheckResponse is the response type for the healthcheck endpoint
type healthCheckResponse struct {
	Environment string `json:"environment"`
	Service     string `json:"service"`
	Version     int    `json:"version"`
}

// healthCheck is a handler used during deployment and monitoring whilst the service runs
func healthCheck(conf config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := healthCheckResponse{
			Environment: conf.Environment.Name,
			Service:     conf.Service.Name,
			Version:     conf.Service.Version,
		}
		utils.WriteResponse(w, r, body, http.StatusOK)
	}
}
