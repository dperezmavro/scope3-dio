package server

import (
	"net/http"

	"github.com/scope3-dio/config"
)

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

// "rows": [
//  {
//  "country": "US",
//  "channel": "web",
//  "impressions": 1000,
//  "inventoryId": "nytimes.com",
//  "utcDatetime": "2024-10-31"
//  }
//  ]
// }
