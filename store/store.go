package store

import (
	"sync"
)

type Store struct {
	values map[string]string
	revision int64
	mu	 sync.RWMutex
	watchers []*watcher
}

func NewStore() *Store {
	// Only Map needs initialization, the rest already have their zero values
	return &Store {
		values: make(map[string]string),
	}
}

func (s *Store) Revision() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.revision
}
	
func (s *Store) Set(key, value string) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.values[key] = value
	s.revision++
	return s.revision
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
		return false, s.revision
	}

	delete(s.values, key)
	s.revision++
	return true, s.revision
}

