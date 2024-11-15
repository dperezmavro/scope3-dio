package storage

import (
	"context"
	"fmt"

	"github.com/scope3-dio/common"
	"github.com/scope3-dio/logging"
)

// generic, we should be able to swap this out
type StorageClient struct {
	initialSize   int
	memoryStorage map[string]string
	remoteFetcher Fetcher
}

func New(initialSize int, f Fetcher) (*StorageClient, error) {

	m := make(map[string]string, initialSize)

	return &StorageClient{
		initialSize:   initialSize,
		memoryStorage: m,
		remoteFetcher: f,
	}, nil
}

func (s *StorageClient) Get(ctx context.Context, queries []common.PropertyQuery) map[string]string {
	res := make(map[string]string, len(queries))
	for _, pq := range queries {
		localSotrageIndex := fmt.Sprintf("%s-%s", pq.UtcDateTime, pq.InventoryID)
		localRes := s.memoryStorage[localSotrageIndex]
		if localRes == "" {
			logging.Info(ctx, logging.Data{"property": localRes}, "property not found locally")
			// TODO: fetch it, use channel here
		} else {
			res[pq.InventoryID] = localRes
		}

	}

	return res
}
