package server

import (
	"context"
	"crypto/subtle"
	"net/http"

	"github.com/google/uuid"
	"github.com/scope3-dio/common"
	"github.com/scope3-dio/logging"
)

func traceIdMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set trace ID on ctx
		traceID := r.Header.Get(common.HeaderTraceID)
		newTraceID := traceID
		if newTraceID == "" {
			newTraceID = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), common.CtxKeyTraceID, newTraceID)
		if traceID == "" {
			logging.Info(ctx, nil, "no X-Trace-Id header, generating new trace ID")
		}
		w.Header().Set(common.HeaderTraceID, newTraceID)

		r = r.WithContext(ctx)
		next(w, r)
	}
}

func authMiddleware(key string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get(common.HeaderAuthorization)
			if key == "" || subtle.ConstantTimeCompare([]byte(authHeader), []byte(key)) != 1 {
				authLog := authHeader
				if len(authHeader) >= 4 {
					authLog = "****" + authHeader[len(authHeader)-4:]
				}
				logging.Info(r.Context(), logging.Data{
					"user-key": authLog,
				}, "unauthorised")

				writeResponse(w, r, map[string]string{"error": "unauthorised"}, http.StatusUnauthorized)
				return
			}
			next(w, r)
		}
	}
}
