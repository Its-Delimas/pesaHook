package store

import (
	"sync"

	"github.com/Its-Delimas/pesaHook/internal/apikey"
)

type MemoryAPIKeyStore struct {
	mu   sync.RWMutex
	data map[string]apikey.APIKey
}

func NewMemoryAPIKeyStore() *MemoryAPIKeyStore {
	return &MemoryAPIKeyStore{data: make(map[string]apikey.APIKey)}
}

func (s *MemoryAPIKeyStore) Save(k apikey.APIKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[k.KeyHash] = k
	return nil
}

func (s *MemoryAPIKeyStore) GetByHash(hash string) (apikey.APIKey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	k, ok := s.data[hash]
	if !ok {
		return apikey.APIKey{}, ErrNotFound
	}
	return k, nil
}
