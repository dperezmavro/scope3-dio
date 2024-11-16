package server

import (
	"net/http"
)

// metrics is a handler for monitoring whilst the service runs
func metrics(sc StorageClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := sc.Metrics(r.Context()).String()
		writeResponse(w, r, body, http.StatusOK)
	}
}
