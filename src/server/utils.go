package server

import (
	"encoding/json"
	"net/http"

	"github.com/scope3-dio/logging"
)

func writeResponse(w http.ResponseWriter, r *http.Request, b interface{}, c int) {
	w.WriteHeader(c)
	w.Header().Set("Content-Type", "application/json")

	resp, err := json.Marshal(b)
	if err != nil {
		logging.Fatal(r.Context(), err, logging.Data{
			"body": b,
		}, "error marshalling response")
	}

	_, err = w.Write(resp)
	if err != nil {
		logging.Error(r.Context(), err, nil, "failed to send response body")
	}
	return
}
