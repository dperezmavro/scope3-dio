package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/scope3-dio/common"
	"github.com/scope3-dio/config"
	"github.com/scope3-dio/logging"
	"github.com/scope3-dio/server"
)

func main() {
	ctx := context.WithValue(context.Background(), common.TraceIdKey, common.BackgroundTraceId)

	conf, err := config.New()
	if err != nil {
		logging.Fatal(ctx, err, nil, "unable to initialise config")
	}

	routerConfig := server.RouterConfig{
		Conf: conf,
	}

	router := server.CreateRouter(routerConfig)

	// start server
	var httpHandler http.Handler = router

	logging.Info(ctx, logging.Data{
		"port": conf.Port,
	}, "starting service")
	err = http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), httpHandler)
	if err != nil {
		logging.Fatal(ctx, err, nil, "service crashed")
	}
}
