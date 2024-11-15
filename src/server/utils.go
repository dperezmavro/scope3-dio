package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/scope3-dio/logging"
)

func writeResponse(b interface{}, c int) http.HandlerFunc {
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

		_, err = w.Write(resp)
		if err != nil {
			logging.Error(r.Context(), err, nil, "failed to send response body")
		}
	}
}
