package server

import (
	"net/http"

	"github.com/dperezmavro/scope3-dio/src/utils"
)

type metricsResponse struct {
	KeysAdded uint64  `json:"keysAdded"`
	Misses    uint64  `json:"misses"`
	Ratio     float64 `json:"ratio"`
}

// metrics is a handler for monitoring whilst the service runs
func metrics(sc StorageClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mr := metricsResponse{
			Misses:    sc.Metrics().Misses(),
			KeysAdded: sc.Metrics().KeysAdded(),
			Ratio:     sc.Metrics().Ratio(),
		}
		utils.WriteResponse(w, r, mr, http.StatusOK)
	}
}
