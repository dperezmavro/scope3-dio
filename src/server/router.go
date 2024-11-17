package server

import (
	"github.com/dperezmavro/scope3-dio/src/config"
	"github.com/dperezmavro/scope3-dio/src/server/middleware"
	"github.com/go-chi/chi"
)

func CreateRouter(
	conf config.Config,
	sc StorageClient,
) chi.Router {
	router := chi.NewRouter()

	router.Post("/v2/measure",
		middleware.TraceIDMiddleware(
			middleware.Performance(
				// this auth is temporary, just does a static key check
				middleware.AuthMiddleware(conf.Scope3APIToken)(
					measure(sc),
				),
			),
		),
	)

	// used for deploying and healthcheck at runtime
	router.Get("/healthcheck",
		middleware.TraceIDMiddleware(
			healthCheck(conf),
		),
	)

	// used for performance monitoring
	router.Get("/metrics",
		middleware.TraceIDMiddleware(
			metrics(sc),
		),
	)

	return router
}
