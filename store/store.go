package store

import (
	"sync"
)

type Store struct {
	values map[string]string
	mu	 sync.RWMutex

}

func NewStore() *Store {
	m := make(map[string]string)
	return &Store{
		values: m,
		mu:     sync.RWMutex{},
	}
}

	
func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.values[key] = value
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.values[key]
	return value, ok
}

func (s *Store) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	// vlaue, ok := s.values[key], ok is true if key exists, false otherwise
	_, ok := s.values[key]
	if !ok {
		return false
	}

	delete(s.values, key)
	return true	
}

