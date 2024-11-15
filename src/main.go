package main

import (
	"context"

	"github.com/scope3-dio/common"
	"github.com/scope3-dio/config"
	"github.com/scope3-dio/logging"
)

func main() {
	ctx := context.WithValue(context.Background(), common.TraceIdKey, common.BackgroundTraceId)

	c, err := config.New()
	if err != nil {
		logging.Fatal(ctx, err, nil, "unable to initialise config")
	}

	logging.Info(ctx, logging.Data{"config": *c}, "started with config")
}
