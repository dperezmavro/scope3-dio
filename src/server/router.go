package server

import (
	"github.com/go-chi/chi"
	"github.com/scope3-dio/config"
)

func CreateRouter(
	conf config.Config,
	sc3 Scope3Client,
	sc StorageClient,
) chi.Router {
	router := chi.NewRouter()

	router.Post("/v2/measure", measure(conf))

	router.Get("/healthcheck", healthCheck(conf))

	return router
}
