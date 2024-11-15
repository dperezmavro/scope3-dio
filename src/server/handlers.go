package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/scope3-dio/config"
	"github.com/scope3-dio/logging"
)

func ResponseWritter(b interface{}, c int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(c)

		resp, err := json.Marshal(b)
		if err != nil {
			logging.Fatal(context.Background(), err, logging.Data{
				"function": "healthcheck",
				"body":     b,
			}, "error marshalling response")
		}

		// common
		_, err = w.Write(resp)
		if err != nil {
			logging.Error(r.Context(), err, nil, "failed to send response body")
		}
	}
}

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
