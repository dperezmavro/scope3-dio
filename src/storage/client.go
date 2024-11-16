package storage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/scope3-dio/src/common"
	"github.com/scope3-dio/src/logging"
)

const defaultCacheTTL = 24 * time.Hour

func New(
	numberOfCounters int64,
	maxCost int64,
	bufferItems int64,
	errors chan error,
	queries chan common.PropertyQuery,
	results chan common.PropertyResponse,
	wg *sync.WaitGroup,
	waitForMissing bool,
) (*Client, error) {

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

	return &Client{
		errors:                errors,
		queries:               queries,
		results:               results,
		wg:                    wg,
		cache:                 cache,
		waitForMissingResults: waitForMissing,
	}, nil
}

func (s *Client) StartListening(ctx context.Context) {
	// wait for the listenForProperties goroutine
	s.wg.Add(1)
	logging.Info(
		ctx,
		logging.Data{
			"function": "client.StartListening",
			"listener": "listenForResults",
			"package":  "storage",
		},
		"listener starting",
	)
	go listenForResults(s)
}

func listenForResults(c *Client) {
	for {
		property := <-c.results
		ctx := context.WithValue(context.Background(), common.CtxKeyTraceID, "listenForResults")
		logging.Info(ctx, logging.Data{"property": property, "weight": property.Weight}, "storing property")
		ok := c.cache.SetWithTTL(property.IndexName(), property, int64(property.Weight), defaultCacheTTL)
		if !ok {
			err := errors.New("unable to set key")
			logging.Error(ctx, err, logging.Data{"key": property.IndexName(), "result": property.Body}, "save error")
			c.errors <- err
		}
	}
}

func (s *Client) Get(ctx context.Context, queries []common.PropertyQuery) []common.PropertyResponse {
	// pre-allocate to avoid resizing
	res := make([]common.PropertyResponse, len(queries))
	notFound := make(map[string]bool, len(queries))
	foundCounter := 0
	notFoundCounter := 0
	for _, pq := range queries {

		localSotrageIndex := pq.IndexName()
		v, found := s.cache.Get(localSotrageIndex)

		if !found {
			notFound[pq.IndexName()] = true
			notFoundCounter++
			logging.Info(ctx, logging.Data{"property": pq.IndexName()}, "property not found locally")
			s.wg.Add(1)
			go func() {
				defer s.wg.Done()
				s.queries <- pq
			}()
		} else {
			logging.Info(ctx, logging.Data{"property": pq.IndexName()}, "property exists locally")
			res[foundCounter] = v
			foundCounter++
		}

	}

	if !s.waitForMissingResults {
		return res[:foundCounter]
	}

	// wait to see if something is fetched
	logging.Info(ctx, logging.Data{"misses": notFoundCounter, "hits": foundCounter}, "waiting for responses")
	timeout := time.After(45 * time.Millisecond)
	tick := time.Tick(5 * time.Millisecond)
	for notFoundCounter > 0 {
		select {
		case <-timeout:
			logging.Info(ctx, logging.Data{"timeout": true}, "search wait")
			return res[:foundCounter]
		case <-tick:
			logging.Info(ctx, logging.Data{"tick": true}, "search wait")
			for k := range notFound {
				val, found := s.cache.Get(k)
				if found {
					logging.Info(ctx, logging.Data{"tick": true, "found": true, "val": val}, "search wait")
					res[foundCounter] = val
					foundCounter++
					notFoundCounter--
					delete(notFound, k)
				}
			}
		}
	}

	// // only return filled slots
	return res[:foundCounter]
}

func (s *Client) Metrics() *ristretto.Metrics {
	return s.cache.Metrics
}
