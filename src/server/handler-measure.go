package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/dperezmavro/scope3-dio/src/common"
	"github.com/dperezmavro/scope3-dio/src/logging"
	"github.com/dperezmavro/scope3-dio/src/utils"
)

// measure is the main query api. uses a local storage client for keeping data.
func measure(sc StorageClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		b, err := io.ReadAll(r.Body)
		if err != nil {
			utils.WriteResponse(w, r, map[string]string{"error": "unable to read body"}, http.StatusInternalServerError)
			logging.Error(ctx, err, logging.Data{"data": string(b)}, "unable to read body")
			return
		}
		defer r.Body.Close()

		var data common.MeasureAPIRequest
		err = json.Unmarshal(b, &data)
		if err != nil {
			utils.WriteResponse(w, r, map[string]string{"error": "unable to unmarshal body"}, http.StatusInternalServerError)
			logging.Error(ctx, err, logging.Data{"data": string(b)}, "unable to unmarshal body")
			return
		}

		if len(data.Rows) < 1 {
			utils.WriteResponse(w, r, map[string]string{"error": "rows must not be empty"}, http.StatusBadRequest)
			logging.Error(ctx, err, logging.Data{"data": string(b)}, "rows must not be empty")
			return
		}

		// do some sanity checking on the values
		validRows := make([]common.PropertyQuery, len(data.Rows))
		validRowsCounter := 0
		for _, row := range data.Rows {
			err := row.Validate(ctx)
			if err != nil {
				logging.Error(
					ctx,
					errors.New("missing val"),
					logging.Data{
						"param":    "UtcDateTime",
						"function": "measure",
					},
					"missing value",
				)
				continue
			}
			validRows[validRowsCounter] = row
			validRowsCounter++
		}

		results := sc.Get(ctx, validRows[:validRowsCounter])
		resp, err := json.Marshal(results)
		if err != nil {
			logging.Error(ctx, err, nil, "marshal error")
		}
		output := resp

		w.Header().Set(common.HeaderContentType, common.HeaderValueContentTypeJSON)
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(output)
		if err != nil {
			logging.Error(ctx, err, nil, "failed to send response body")
		}
	}
}
