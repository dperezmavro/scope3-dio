package storage

import (
	"context"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/dperezmavro/scope3-dio/src/common"
	"github.com/dperezmavro/scope3-dio/src/logging"
)

func NewStorageImplementation(
	numberOfCounters int64,
	maxCost int64,
	bufferItems int64,
) (*Ristretto, error) {

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

	return &Ristretto{
		r: cache,
	}, nil
}

type Ristretto struct {
	r *ristretto.Cache[string, common.PropertyResponse]
}

func (r *Ristretto) Metrics() *ristretto.Metrics {
	return r.r.Metrics
}

func (r *Ristretto) Get(key string) (common.PropertyResponse, bool) {
	return r.r.Get(key)
}

func (r *Ristretto) SetWithTTL(key string, value common.PropertyResponse, cost int64, ttl time.Duration) bool {
	return r.r.SetWithTTL(key, value, cost, ttl)
}
