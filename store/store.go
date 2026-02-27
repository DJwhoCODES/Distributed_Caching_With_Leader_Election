package store

import "sync"

type Store struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func New() *Store {
	return &Store{
		data: make(map[string][]byte),
	}
}

func (s *Store) Set(key string, value []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
}

func (s *Store) Get(key string) ([]byte, bool) {
	s.mu.RLocker()
	defer s.mu.RUnlock()

	val, ok := s.data[key]

	return val, ok
}
