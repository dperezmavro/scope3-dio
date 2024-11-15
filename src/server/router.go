package server

import (
	"github.com/go-chi/chi"
)

func CreateRouter(rc RouterConfig) chi.Router {
	router := chi.NewRouter()

	router.Get("/healthcheck", healthCheckHandler(*rc.Conf))

	return router
}
