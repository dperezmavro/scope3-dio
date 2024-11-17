package storage

import (
	"sync"
	"time"

	"github.com/dperezmavro/scope3-dio/src/common"
)

// Client is an is responsible for managing our storage.
type Client struct {
	// channels for communicating with other goroutes to share results
	errors  chan error
	queries chan common.PropertyQuery
	results chan common.PropertyResponse
	wg      *sync.WaitGroup

	// in-memory cache
	cache                 Implementation
	waitForMissingResults bool
}

type Implementation interface {
	Get(key string) (common.PropertyResponse, bool)
	SetWithTTL(key string, value common.PropertyResponse, cost int64, ttl time.Duration) bool
}
