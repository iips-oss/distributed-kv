package store

import "sync"

type Store struct {
	data map[string]string
	mu   sync.RWMutex
	wal  *WAL
}

func NewStore(wal *WAL) *Store {
	return &Store{
		data: make(map[string]string),
		wal:  wal,
	}
}

func (s *Store) Put(key, value string) error {
	// 1. we will write the operation to the WAL before applying it to memoryso that if the server crashes after writing to the WAL but before updating memory we can still recover the operation when we replay the WAL on startup
	err := s.wal.Write(Entry{
		Op:    "PUT",
		Key:   key,
		Value: value,
	})
	if err != nil {
		return err
	}

	// 2. here we will  update memory
	s.mu.Lock()
	s.data[key] = value
	s.mu.Unlock()

	return nil
}
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.data[key]
	return val, ok
}
func (s *Store) Delete(key string) error {
	// 1. we will write the operation to the WAL before applying it to memory so that if the server crashes after writing to the WAL but before updating memory we can still recover the operation when we replay the WAL on startup
	err := s.wal.Write(Entry{
		Op:  "DELETE",
		Key: key,
	})
	if err != nil {
		return err
	}

	// 2. here we will  update memory
	s.mu.Lock()
	delete(s.data, key)
	s.mu.Unlock()

	return nil
}

// replay the WAL by reading all entries from the file and applying them to the in-memory store so that we can recover the state of the store after a crash
func (s *Store) Replay() error {
	return s.wal.Replay(func(e Entry) {

		switch e.Op {
		case "PUT":
			s.data[e.Key] = e.Value

		case "DELETE":
			delete(s.data, e.Key)
		}
	})
}
