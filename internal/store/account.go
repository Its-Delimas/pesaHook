package store

import (
	"sync"

	"github.com/Its-Delimas/pesaHook/internal/account"
)

type AccountStore interface {
	Save(a account.Account) error
	GetById(id string) (account.Account, error)
}

type MemoryAccountStore struct {
	mu   sync.RWMutex
	data map[string]account.Account
}

func NewMemoryAccountStore() *MemoryAccountStore {
	return &MemoryAccountStore{data: make(map[string]account.Account)}
}

func (s *MemoryAccountStore) Save(a account.Account) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[a.ID] = a
	return nil
}

func (s *MemoryAccountStore) GetById(id string) (account.Account, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.data[id]
	if !ok {
		return account.Account{}, ErrNotFound
	}
	return a, nil
}
