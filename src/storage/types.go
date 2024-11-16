package storage

import (
	"sync"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/scope3-dio/src/common"
)

// Client is an is responsible for managing our storage.
type Client struct {
	// channels for communicating with other goroutes to share results
	errors  chan error
	queries chan common.PropertyQuery
	results chan common.PropertyResponse
	wg      *sync.WaitGroup

	// in-memory cache
	cache                 *ristretto.Cache[string, common.PropertyResponse]
	waitForMissingResults bool
}
