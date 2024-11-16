package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/scope3-dio/src/clients/scope3"
	"github.com/scope3-dio/src/common"
	"github.com/scope3-dio/src/config"
	"github.com/scope3-dio/src/logging"
	"github.com/scope3-dio/src/server"
	"github.com/scope3-dio/src/storage"
)

var (
	queryChannel    = make(chan common.PropertyQuery)
	responseChannel = make(chan common.PropertyResponse)
	errorChannel    = make(chan error)
	wg              = &sync.WaitGroup{}
)

func main() {
	ctx := context.WithValue(context.Background(), common.CtxKeyTraceID, common.BackgroundTraceID)

	conf, err := config.New()
	if err != nil {
		logging.Fatal(ctx, err, nil, "unable to initialise config")
	}

	scope3Client := scope3.New(
		conf.Scope3APIToken,
		errorChannel,
		queryChannel,
		responseChannel,
		wg,
	)

	// start listening for async fetches
	scope3Client.StartListening(ctx)

	// defaultSize := 10
	storageClient, err := storage.New(
		1e7, 1<<30, 64,
		errorChannel,
		queryChannel,
		responseChannel,
		wg,
	)
	if err != nil {
		logging.Fatal(ctx, err, nil, "unable to initialise storage client")
	}

	// start listening for async stores
	storageClient.StartListening(ctx)

	// start server
	var httpHandler http.Handler = server.CreateRouter(*conf, storageClient)

	logging.Info(ctx, logging.Data{
		"port": conf.Port,
	}, "starting service")
	err = http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), httpHandler)
	if err != nil {
		logging.Fatal(ctx, err, nil, "service crashed")
	}

	wg.Wait()
}
