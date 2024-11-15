package storage

// generic, we should be able to swap this out
type StorageClient struct {
	initialSize   int
	memoryStorage map[string]string
}

func New(initialSize int) (*StorageClient, error) {

	m := make(map[string]string, initialSize)

	return &StorageClient{
		initialSize:   initialSize,
		memoryStorage: m,
	}, nil
}

func (s *StorageClient) Get(k string) string {
	return s.memoryStorage[k]
}
