package server

import (
	"context"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/dperezmavro/scope3-dio/src/common"
)

type StorageClient interface {
	Get(context.Context, []common.PropertyQuery) []common.PropertyResponse
	Metrics() *ristretto.Metrics
}
