package server

import (
	"context"
	"net/http"

	"github.com/scope3-dio/config"
)

type StorageClient interface {
	Get(context.Context, string) string
}

type RouterConfig struct {
	Conf        *config.Config
	HealthCheck http.HandlerFunc
}

type healthCheckResponse struct {
	Environment string `json:"environment"`
	Service     string `json:"service"`
	Version     int    `json:"version"`
}
