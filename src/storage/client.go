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
	errors  chan error
	queries chan common.PropertyQuery
	results chan common.PropertyResponse
	wg      *sync.WaitGroup

	cache *ristretto.Cache[string, string]
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

	cache, err := ristretto.NewCache(&ristretto.Config[string, string]{
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
		errors <- err
		return nil, err
	}

	// // set a value with a cost of 1
	// cache.Set("key", "value", 1)

	// // wait for value to pass through buffers
	// cache.Wait()

	// // get value from cache
	// value, found := cache.Get("key")

	// logging.Info(ctx, logging.Data{"properties": properties, "function": "listenForProperties"}, "storing property")

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
	logging.Info(ctx, logging.Data{"function": "client.StartListening", "listener": "listenForResults"}, "starting storage goroutine listener")
	go listenForResults(
		s,
		s.results,
		s.errors,
	)
}

func listenForResults(
	c *Client,
	results chan common.PropertyResponse,
	errChan chan error,
) {
	for {
		property := <-results
		ctx := context.WithValue(context.Background(), common.CtxKeyTraceID, "listenForResults")
		logging.Info(ctx, logging.Data{"property": property}, "saving property")
		ok := c.cache.Set(property.IndexName(), property.Body, property.Weight)
		if !ok {
			err := errors.New("unable to set key")
			logging.Error(ctx, err, logging.Data{"key": property.IndexName(), "result": property.Body}, "save error")
			errChan <- err
		}
	}
}

func (s *Client) Get(ctx context.Context, queries []common.PropertyQuery) []string {
	res := make([]string, len(queries))
	for _, pq := range queries {

		localSotrageIndex := pq.IndexName()
		localRes := ""
		v, found := s.cache.Get(localSotrageIndex)

		if !found {
			logging.Info(ctx, logging.Data{"property": localRes}, "property not found locally")
			go func() {
				s.wg.Add(1)
				defer s.wg.Done()
				s.queries <- pq
			}()
		} else {
			logging.Info(ctx, logging.Data{"property": localRes}, "property exists locally")
			res = append(res, v)
		}

	}

	return res
}
