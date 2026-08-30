package store

import (
	"sync"
	"time"

	"github.com/Its-Delimas/pesahook/internal/event"
)

type DeadLetter struct {
	Event      event.NormalizedEvent
	EndpointID string
	LastError  string
	FailedAt   time.Time
	Attempts   int
}

type DeadLetterStore interface {
	Save(dl DeadLetter) error
	ListByEndpoint(endpointID string) ([]DeadLetter, error)
}

type MemoryDeadLetterStore struct {
	mu   sync.RWMutex
	data []DeadLetter
}

func NewMemoryDeadLetterStore() *MemoryDeadLetterStore {
	return &MemoryDeadLetterStore{}
}

func (s *MemoryDeadLetterStore) Save(dl DeadLetter) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = append(s.data, dl)
	return nil
}

func (s *MemoryDeadLetterStore) ListByEndpoint(endpointID string) ([]DeadLetter, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var results []DeadLetter
	for _, dl := range s.data {
		if dl.EndpointID == endpointID {
			results = append(results, dl)
		}
	}
	return results, nil
}
