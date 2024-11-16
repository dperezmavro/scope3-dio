package server

import (
	"github.com/go-chi/chi"
	"github.com/scope3-dio/src/config"
)

func CreateRouter(
	conf config.Config,
	sc StorageClient,
) chi.Router {
	router := chi.NewRouter()

	router.Post("/v2/measure",
		traceIDMiddleware(
			// this auth is temporary, just does a static key check
			authMiddleware(conf.Scope3APIToken)(
				measure(sc),
			),
		),
	)

	router.Get("/healthcheck",
		traceIDMiddleware(
			healthCheck(conf),
		),
	)

	router.Get("/metrics",
		traceIDMiddleware(
			metrics(sc),
		),
	)

	return router
}
