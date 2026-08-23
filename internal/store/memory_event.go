package store

import (
	"sync"

	"github.com/Its-Delimas/pesaHook/internal/event"
)

type MemoryEventStore struct {
	mu   sync.RWMutex
	data map[string]event.NormalizedEvent
}

func NewMemoryEVentStore() *MemoryEventStore {
	return &MemoryEventStore{data: make(map[string]event.NormalizedEvent)}
}

func (s *MemoryEventStore) Save(e event.NormalizedEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[e.EventID] = e
	return nil
}

func (s *MemoryEventStore) GetByID(id string) (event.NormalizedEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.data[id]
	if !ok {
		return event.NormalizedEvent{}, ErrNotFound
	}
	return e, nil
}

func (s *MemoryEventStore) ListByEndpoint(endpointID string) ([]event.NormalizedEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var results []event.NormalizedEvent
	for _, e := range s.data {
		if e.ProviderMeta["endpoint_id"] == endpointID {
			results = append(results, e)
		}
	}
	return results, nil
}
