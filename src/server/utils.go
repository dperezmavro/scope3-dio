package server

import (
	"encoding/json"
	"net/http"

	"github.com/scope3-dio/src/common"
	"github.com/scope3-dio/src/logging"
)

func writeResponse(w http.ResponseWriter, r *http.Request, b interface{}, c int) {
	w.WriteHeader(c)
	w.Header().Set(common.HeaderContentType, common.HeaderValueContentTypeJson)

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
}
