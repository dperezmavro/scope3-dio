package middleware

import (
	"context"
	"net/http"

	"github.com/dperezmavro/scope3-dio/src/common"
	"github.com/dperezmavro/scope3-dio/src/logging"
	"github.com/google/uuid"
)

func TraceIDMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set trace ID on ctx
		traceID := r.Header.Get(common.HeaderTraceID)
		newTraceID := traceID
		if newTraceID == "" {
			newTraceID = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), common.CtxKeyTraceID, newTraceID)
		if traceID == "" {
			logging.Info(ctx, logging.Data{"traceId": newTraceID}, "no X-Trace-Id header, generating new trace ID")
		}
		w.Header().Set(common.HeaderTraceID, newTraceID)

		r = r.WithContext(ctx)
		next(w, r)
	}
}
