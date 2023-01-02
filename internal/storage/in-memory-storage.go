package storage

import (
	"fmt"
	"net/url"
)

type InMemoryStorage struct {
	key2url map[string]*url.URL
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		key2url: make(map[string]*url.URL),
	}
}

func (s *InMemoryStorage) Put(key string, url *url.URL) (err error) {
	if _, isOk := s.Get(key); isOk {
		return fmt.Errorf("key %q already exists", key)
	}
	s.key2url[key] = url
	return
}

func (s *InMemoryStorage) Get(key string) (url *url.URL, isOk bool) {
	url, isOk = s.key2url[key]
	return
}
