package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/scope3-dio/clients/scope3"
	"github.com/scope3-dio/common"
	"github.com/scope3-dio/config"
	"github.com/scope3-dio/logging"
	"github.com/scope3-dio/server"
	"github.com/scope3-dio/storage"
)

func main() {
	ctx := context.WithValue(context.Background(), common.TraceIdKey, common.BackgroundTraceId)

	conf, err := config.New()
	if err != nil {
		logging.Fatal(ctx, err, nil, "unable to initialise config")
	}

	scope3Client := scope3.New(conf.Scope3APIToken)
	defaultSize := 10
	storageClient, err := storage.New(defaultSize)
	if err != nil {
		logging.Fatal(ctx, err, logging.Data{"size": defaultSize}, "unable to initialise storage client")
	}

	// start server
	var httpHandler http.Handler = server.CreateRouter(*conf, scope3Client, storageClient)

	logging.Info(ctx, logging.Data{
		"port": conf.Port,
	}, "starting service")
	err = http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), httpHandler)
	if err != nil {
		logging.Fatal(ctx, err, nil, "service crashed")
	}
}
