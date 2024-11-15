package server

import (
	"github.com/go-chi/chi"
	"github.com/scope3-dio/config"
)

func CreateRouter(
	conf config.Config,
	sc StorageClient,
) chi.Router {
	router := chi.NewRouter()

	router.Post("/v2/measure",
		traceIdMiddleware(
			// this auth is temporary, just does a static key check
			authMiddleware(conf.Scope3APIToken)(
				measure(sc),
			),
		),
	)

	router.Get("/healthcheck", healthCheck(conf))

	return router
}
