package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/dperezmavro/scope3-dio/src/common"
	"github.com/dperezmavro/scope3-dio/src/logging"
	"github.com/dperezmavro/scope3-dio/src/utils"
)

func AuthMiddleware(key string) func(http.HandlerFunc) http.HandlerFunc {
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

				utils.WriteResponse(w, r, map[string]string{"error": "unauthorised"}, http.StatusUnauthorized)
				return
			}
			next(w, r)
		}
	}
}
