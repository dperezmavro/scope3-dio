package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/scope3-dio/src/common"
	"github.com/scope3-dio/src/logging"
)

// measure is the main query api. uses a local storage client for keeping data.
func measure(sc StorageClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		b, err := io.ReadAll(r.Body)
		if err != nil {
			writeResponse(w, r, map[string]string{"error": "unable to read body"}, http.StatusInternalServerError)
			logging.Error(ctx, err, logging.Data{"data": string(b)}, "unable to read body")
			return
		}
		defer r.Body.Close()

		var data common.MeasureAPIRequest
		err = json.Unmarshal(b, &data)
		if err != nil {
			writeResponse(w, r, map[string]string{"error": "unable to unmarshal body"}, http.StatusInternalServerError)
			logging.Error(ctx, err, logging.Data{"data": string(b)}, "unable to unmarshal body")
			return
		}

		// do some sanity checking on the values
		for _, row := range data.Rows {
			if row.InventoryID == "" {
				logging.Error(ctx, errors.New("missing val"), logging.Data{"param": "InventoryID", "function": "measure"}, "missing value")
				writeResponse(w, r, map[string]string{"error": "missing inventoryId"}, http.StatusBadRequest)
				return
			}

			if row.UtcDateTime == "UtcDateTime" {
				logging.Error(ctx, errors.New("missing val"), logging.Data{"param": "UtcDateTime", "function": "measure"}, "missing value")
				writeResponse(w, r, map[string]string{"error": "missing utcDateTime"}, http.StatusBadRequest)
				return
			}

			if row.Impressions == 0 {
				row.Impressions = 1000
			}

			if row.Weight == 0 {
				row.Weight = 1
			}
		}

		results := sc.Get(ctx, data.Rows)
		output := []byte(fmt.Sprintf(`{"rows": [%s]}`, strings.Join(results, ",")))

		w.Header().Set(common.HeaderContentType, common.HeaderValueContentTypeJSON)
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(output)
		if err != nil {
			logging.Error(ctx, err, nil, "failed to send response body")
		}
	}
}
