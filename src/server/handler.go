package server

import (
	"encoding/json"
	"net/http"

	"github.com/scope3-dio/config"
	"github.com/scope3-dio/logging"
)

func healthCheckHandler(conf config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := healthCheckResponse{
			Environment: conf.Environment.Name,
			Service:     conf.Service.Name,
			Version:     conf.Service.Version,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp, err := json.Marshal(body)
		if err != nil {
			logging.Fatal(r.Context(), err, logging.Data{
				"function": "healthcheck",
				"body":     body,
			}, "error marshalling response")
		}

		_, err = w.Write(resp)
		if err != nil {
			logging.Error(r.Context(), err, nil, "failed to send response body")
		}
	}
}
