package storage

import (
	"context"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/dperezmavro/scope3-dio/src/common"
	"github.com/dperezmavro/scope3-dio/src/logging"
)

func NewStorageImplementation(
	numberOfCounters int64,
	maxCost int64,
	bufferItems int64,
) (*ristretto.Cache[string, common.PropertyResponse], error) {

	cache, err := ristretto.NewCache(&ristretto.Config[string, common.PropertyResponse]{
		NumCounters: numberOfCounters,
		MaxCost:     maxCost,
		BufferItems: bufferItems,
	})
	if err != nil {
		logging.Error(context.Background(), err, logging.Data{
			"NumCounters": numberOfCounters,
			"MaxCost":     maxCost,
			"BufferItems": bufferItems,
		}, "error in NewCache")
		return nil, err
	}

	return cache, nil
}
