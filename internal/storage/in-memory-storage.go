package storage

import (
	"github.com/teris-io/shortid"
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

func (s *InMemoryStorage) Put(key string, url *url.URL) (retKey string, err error) {
	if key == "" {
		key, err = s.GenerateNewKey()
		if err != nil {
			return
		}
	}
	s.key2url[key] = url
	retKey = key
	return
}

func (s *InMemoryStorage) Get(key string) (url *url.URL, isOk bool) {
	url, isOk = s.key2url[key]
	return
}

func (s *InMemoryStorage) GenerateNewKey() (ret string, err error) {
	for {
		ret, err = shortid.Generate()
		if err != nil {
			return
		}
		if _, isOk := s.Get(ret); isOk {
			continue
		}
		break
	}

	return
}
