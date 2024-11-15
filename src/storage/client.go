package storage

import (
	"context"

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

func (s *StorageClient) Get(ctx context.Context, k string) string {
	res := s.memoryStorage[k]

	if res == "" {
		logging.Info(ctx, logging.Data{"property": k}, "property not found locally")
	}

	return res
}
