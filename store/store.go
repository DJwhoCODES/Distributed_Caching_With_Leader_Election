package store

import (
	"sync"
	"time"
)

type item struct {
	value      []byte
	expiration int64
}

type Store struct {
	mu   sync.RWMutex
	data map[string]item
}

func New() *Store {
	s := &Store{
		data: make(map[string]item),
	}
	go s.cleanupLoop()
	return s
}

func (s *Store) Set(key string, value []byte, ttl int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(time.Duration(ttl) * time.Millisecond).UnixMilli()
	}

	s.data[key] = item{
		value:      value,
		expiration: exp,
	}
}

func (s *Store) Get(key string) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	it, ok := s.data[key]

	if !ok {
		return nil, false
	}

	if it.expiration > 0 && time.Now().UnixMilli() > it.expiration {
		s.mu.Lock()
		delete(s.data, key)
		s.mu.Unlock()
		return nil, false
	}

	return it.value, true
}

func (s *Store) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		now := time.Now().UnixMilli()

		s.mu.Lock()
		for k, it := range s.data {
			if it.expiration > 0 && now > it.expiration {
				delete(s.data, k)
			}
		}
		s.mu.Unlock()
	}
}
