package store

import (
	"sync"
)

type Store struct {
	values map[string]string
	mu	 sync.RWMutex
	Revision int64
}

func NewStore() *Store {
	m := make(map[string]string)
	return &Store{
		values: m,
		mu:     sync.RWMutex{},
		Revision: 0,
	}
}

	
func (s *Store) Set(key, value string) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.values[key] = value
	s.Revision++
	return s.Revision
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.values[key]
	return value, ok
}

func (s *Store) Delete(key string) (bool, int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// vlaue, ok := s.values[key], ok is true if key exists, false otherwise
	_, ok := s.values[key]
	if !ok {
		return false, s.Revision
	}

	delete(s.values, key)
	s.Revision++
	return true, s.Revision
}

