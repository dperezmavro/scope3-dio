package middleware

import (
	"net/http"
	"time"

	"github.com/dperezmavro/scope3-dio/src/logging"
)

func Performance(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		next(w, r)
		logging.Info(
			r.Context(),
			logging.Data{
				"unit":     "microseconds",
				"duration": time.Since(now).Microseconds(),
			},
			"duration",
		)
	}
}
