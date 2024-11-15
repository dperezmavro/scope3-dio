package main

import (
	"context"
	"log"

	"github.com/scope3-dio/config"
	"github.com/scope3-dio/logging"
)

func main() {
	ctx := context.Background()

	c, err := config.New()
	if err != nil {
		logging.Fatal(ctx, err, nil, "unable to initialise config")
	}

	log.Printf("started with config %+v", c)
}
