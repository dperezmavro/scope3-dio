package server

import (
	"errors"
	"net/http"

	"github.com/scope3-dio/logging"
)

const (
	paramChannel     = "channel"
	paramCountry     = "country"
	paramInventoryId = "inventoryId"
	paramUtcDateTime = "utcDatetime"
)

func measure(sc StorageClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logging.Info(r.Context(), nil, "in measure")
		values := r.PostForm

		if values.Get(paramInventoryId) == "" {
			logging.Error(r.Context(), errors.New("missing val"), nil, "in measure")
			writeResponse(w, r, map[string]string{"error": "missing paramInventoryId"}, http.StatusBadRequest)
			return
		}

		if values.Get(paramUtcDateTime) == "" {
			writeResponse(w, r, map[string]string{"error": "missing paramUtcDateTime"}, http.StatusBadRequest)
			return
		}

		res := sc.Get(ctx, "abc")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(res))
		if err != nil {
			logging.Error(r.Context(), err, nil, "failed to send response body")
		}
	}
}
