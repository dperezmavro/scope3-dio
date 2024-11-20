package storage

import (
	"context"
	"errors"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/dperezmavro/scope3-dio/src/common"
	"github.com/dperezmavro/scope3-dio/src/logging"
)

const defaultCacheTTL = 24 * time.Hour

func New(
	cache Implementation,
	errorChan chan error,
	queries chan []common.PropertyQuery,
	results chan []common.PropertyResponse,
	waitForMissing bool,
) (*Client, error) {

	return &Client{
		errors:                errorChan,
		queries:               queries,
		results:               results,
		cache:                 cache,
		waitForMissingResults: waitForMissing,
	}, nil
}

func (s *Client) StartListening(ctx context.Context) {
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
		properties := <-c.results
		ctx := context.WithValue(context.Background(), common.CtxKeyTraceID, "listenForResults")
		logging.Info(ctx, logging.Data{"properties": properties}, "storing properties")
		for _, pr := range properties {
			ok := c.cache.SetWithTTL(pr.IndexName(), pr, int64(pr.Weight), defaultCacheTTL)
			if !ok {
				err := errors.New("unable to set key")
				logging.Error(ctx, err, logging.Data{"key": pr.IndexName(), "result": pr.TotalEmissions}, "save error")
				c.errors <- err
			}
		}
	}
}

func (s *Client) Get(ctx context.Context, queries []common.PropertyQuery) []common.PropertyResponse {
	// pre-allocate to avoid resizing
	res := make([]common.PropertyResponse, len(queries))
	notFoundMap := make(map[string]bool, len(queries))
	notFoundSlice := make([]common.PropertyQuery, len(queries))
	foundCounter := 0
	notFoundCounter := 0
	for _, pq := range queries {
		localSotrageIndex := pq.IndexName()
		v, found := s.cache.Get(localSotrageIndex)

		if !found {
			notFoundMap[pq.IndexName()] = true
			notFoundSlice[notFoundCounter] = pq
			notFoundCounter++
			logging.Info(ctx, logging.Data{"property": pq.IndexName()}, "property not found locally")
		} else {
			logging.Info(ctx, logging.Data{"property": pq.IndexName()}, "property exists locally")
			res[foundCounter] = v
			foundCounter++
		}

	}

	go func() {
		s.queries <- notFoundSlice[:notFoundCounter]
	}()

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
			for k := range notFoundMap {
				val, found := s.cache.Get(k)
				if found {
					logging.Info(ctx, logging.Data{"tick": true, "found": true, "val": val}, "search wait")
					res[foundCounter] = val
					foundCounter++
					notFoundCounter--
					delete(notFoundMap, k)
				}
			}
		}
	}

	// // only return filled slots
	return res[:foundCounter]
}

func (s *Client) Metrics() *ristretto.Metrics {
	return s.cache.Metrics()
}
