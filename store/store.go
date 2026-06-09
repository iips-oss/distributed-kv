package store

import (
	"strings"
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
	
// We are breaking the function in two parts, as we don't want to hold while notify runs
func (s *Store) set(key, value string) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = value
	s.revision++
	return s.revision
}

func (s *Store) Set(key, value string) int64 {
	rev := s.set(key, value)
	s.notify(key, value, SET, rev)
	return rev
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.values[key]
	return value, ok
}

func (s *Store) delete(key string) (bool, int64) {
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

func (s *Store) Delete(key string) (bool, int64) {
	deleted, rev := s.delete(key)
	if deleted {
		s.notify(key, "", DELETE, rev)
		return deleted, rev
	}
	return deleted, rev
}

// To notify events for changes in the watchers prefix, through it's channel
func (s *Store) notify(key string, value string, command Command, revision int64) {
	// for index, value := range myMap {...}
	for _, w := range s.watchers{
		if (strings.HasPrefix(key, w.Prefix)) {
			w.events <- Event{
				Command: command,
				Key: key,
				Value: value,
				Revision: revision,
			}
		}
	}
}

