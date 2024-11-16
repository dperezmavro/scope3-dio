package server

import (
	"net/http"
)

type metricsResponse struct {
	KeysAdded uint64  `json:"keysAdded"`
	Misses    uint64  `json:"misses"`
	Ratio     float64 `json:"ratio"`
}

// metrics is a handler for monitoring whilst the service runs
func metrics(sc StorageClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := sc.Metrics()
		mr := metricsResponse{
			Misses:    m.Misses(),
			KeysAdded: m.KeysAdded(),
			Ratio:     m.Ratio(),
		}
		writeResponse(w, r, mr, http.StatusOK)
	}
}
