package server

import (
	"github.com/go-chi/chi"
	"github.com/scope3-dio/config"
)

func CreateRouter(conf config.Config) chi.Router {
	router := chi.NewRouter()

	router.Get("/v2/measure", measure(conf))

	router.Get("/healthcheck", healthCheck(conf))

	return router
}
