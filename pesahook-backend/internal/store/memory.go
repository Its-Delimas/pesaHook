package store

import (
	"sync"

	"github.com/Its-Delimas/pesahook/internal/endpoint"
)

type MemoryEndpointStore struct {
	mu   sync.RWMutex
	data map[string]endpoint.Endpoint
}

func NewMemoryEndpointStore() *MemoryEndpointStore {
	return &MemoryEndpointStore{data: make(map[string]endpoint.Endpoint)}
}

func (s *MemoryEndpointStore) ListByAccount(accountID string) ([]endpoint.Endpoint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var results []endpoint.Endpoint
	for _, e := range s.data {
		if e.AccountID == accountID {
			results = append(results, e)
		}
	}
	return results, nil
}

func (s *MemoryEndpointStore) Save(e endpoint.Endpoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[e.ID] = e
	return nil
}

func (s *MemoryEndpointStore) GetByID(id string) (endpoint.Endpoint, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.data[id]
	if !ok {
		return endpoint.Endpoint{}, ErrNotFound
	}
	return e, nil
}
