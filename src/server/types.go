package server

import (
	"net/http"

	"github.com/scope3-dio/config"
)

type Scope3Client interface {
	FetchProperty(string) error
}

type StorageClient interface {
	Get(string) string
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

type MeasureApiRequest struct {
	Rows []RowItem `json:"rows"`
}

type RowItem struct {
	Country     string `json:"country"`
	Channel     string `json:"channel"`
	Impressions int    `json:"impressions"`
	InventoryID string `json:"InventoryID"`
	UtcDateTime string `json:"utcDatetime"`
}
