package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/scope3-dio/src/common"
	"github.com/scope3-dio/src/logging"
)

// Client is an is responsible for managing our storage.
type Client struct {
	// channels for communicating with other goroutes to share results
	errors  chan error
	queries chan common.PropertyQuery
	results chan common.PropertyResponse
	wg      *sync.WaitGroup

	// in-memory cache
	cache *ristretto.Cache[string, common.PropertyResponse]
}

func New(
	numberOfCounters int64,
	maxCost int64,
	bufferItems int64,
	errors chan error,
	queries chan common.PropertyQuery,
	results chan common.PropertyResponse,
	wg *sync.WaitGroup,
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
		errors:  errors,
		queries: queries,
		results: results,
		wg:      wg,
		cache:   cache,
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
		ok := c.cache.Set(property.IndexName(), property, int64(property.Weight))
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
	foundCounter := 0
	for _, pq := range queries {

		localSotrageIndex := pq.IndexName()
		v, found := s.cache.Get(localSotrageIndex)

		if !found {
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

	// only return filled slots
	return res[:foundCounter]
}

func (s *Client) Metrics() *ristretto.Metrics {
	return s.cache.Metrics
}
