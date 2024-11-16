package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/dperezmavro/scope3-dio/src/clients/scope3"
	"github.com/dperezmavro/scope3-dio/src/common"
	"github.com/dperezmavro/scope3-dio/src/config"
	"github.com/dperezmavro/scope3-dio/src/logging"
	"github.com/dperezmavro/scope3-dio/src/server"
	"github.com/dperezmavro/scope3-dio/src/storage"
)

var (
	queryChannel    = make(chan common.PropertyQuery)
	responseChannel = make(chan common.PropertyResponse)
	errorChannel    = make(chan error)
	wg              = &sync.WaitGroup{}
)

func main() {
	ctx := context.WithValue(
		context.Background(),
		common.CtxKeyTraceID,
		"initialising",
	)

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

	storageClient, err := storage.New(
		1e7, 1<<30, 64,
		errorChannel,
		queryChannel,
		responseChannel,
		wg,
		conf.WaitForMissing,
	)
	if err != nil {
		logging.Fatal(ctx, err, nil, "unable to initialise storage client")
	}

	// start listening for async store requests
	storageClient.StartListening(ctx)

	// start listening to error channel
	wg.Add(1)
	go func() {
		logging.Info(ctx, logging.Data{"goroutine": "main error listener"}, "listener starting")
		for {
			asyncErr := <-errorChannel
			if asyncErr != nil {
				logging.Error(ctx, asyncErr, logging.Data{"goroutine": "main error listener"}, "error")
			}
		}
	}()

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
