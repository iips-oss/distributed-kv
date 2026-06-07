package controller

import (
	"errors" // to create custom errors
	"sync"   // for synchronization
	"time"   // for TTL
)

// For key-value pairs. Each key stores and item
type item struct {
	value     string    // The actual value of the key
	expiresAt time.Time // TTL (time to live)
}

// Actual db for storing KV pairs
type Store struct {
	data map[string]item
	mu   sync.RWMutex // for lock (to prevent race condition)
}

// Creates a new store
func NewStore() *Store {
	return &Store{
		data: make(map[string]item), // will create a store and return its address
	}
}

// Set method that belongs to the Store structure
// Used for setting / inserting a KV pair in the store (db)
func (s *Store) Set(key, value string, ttlSeconds int) {
	s.mu.Lock()                                                           // Acquiring the lock
	defer s.mu.Unlock()                                                   // Release the lock after function returns
	expiration := time.Now().Add(time.Duration(ttlSeconds) * time.Second) // Setting the TTL
	s.data[key] = item{value: value, expiresAt: expiration}               // setting the data in the store
}

// To retreive the value of a key
// returns the value or error if the key does not exists in the store
func (s *Store) Get(key string) (string, error) {
	s.mu.RLock() // Acquiring a shared read lock
	// Allows multiple readers to read a shared resource at the same time
	// but blocks any writes until all the read locks are released
	defer s.mu.RUnlock() // Release the lock later

	it, exists := s.data[key]

	// If there is no key present or the TTL has expired,
	// delete the key from the store
	if !exists || time.Now().After(it.expiresAt) {
		// Key not found or expired, trying to delete it
		s.mu.RUnlock()      // release the read lock
		s.mu.Lock()         // acquire write lock
		delete(s.data, key) // remove the key from store
		s.mu.Unlock()       // release write lock
		s.mu.RLock()        // acquire read lock again
		return "", errors.New("Key not found or expired")
	}

	return it.value, nil // if key is present then return the value
}

func (s *Store) Delete(key string) error {
	s.mu.Lock()         // acquire write lock
	defer s.mu.Unlock() // lock should not be released until the completion

	// If the key does not exist, throw an error
	if _, exists := s.data[key]; !exists {
		return errors.New("Key not found")
	}

	// Remove the key from the DB if it exists
	delete(s.data, key)
	return nil
}
