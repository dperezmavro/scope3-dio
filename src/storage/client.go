package storage

import (
	"context"
	"sync"

	"github.com/scope3-dio/src/common"
	"github.com/scope3-dio/src/logging"
)

// Client is an is responsible for managing our storage.
type Client struct {
	errors  chan error
	queries chan common.PropertyQuery
	results chan []common.PropertyResponse
	wg      *sync.WaitGroup
}

func New(
	errors chan error,
	queries chan common.PropertyQuery,
	results chan []common.PropertyResponse,
	wg *sync.WaitGroup,
) *Client {

	return &Client{
		errors:  errors,
		queries: queries,
		results: results,
		wg:      wg,
	}
}

func (s *Client) Get(ctx context.Context, queries []common.PropertyQuery) map[string]string {
	res := make(map[string]string, len(queries))
	for _, pq := range queries {

		// 	localSotrageIndex := fmt.Sprintf("%s-%s", pq.UtcDateTime, pq.InventoryID)
		// 	localRes := s.memoryStorage[localSotrageIndex]
		localRes := ""

		if localRes == "" {
			logging.Info(ctx, logging.Data{"property": localRes}, "property not found locally")
			// TODO: fetch it, use channel here
			s.queries <- pq
		} else {
			res[pq.InventoryID] = localRes
		}

	}

	return res
}
