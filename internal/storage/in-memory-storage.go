package storage

import (
	"net/url"
)

type InMemoryStorage struct {
	shid2url map[string]*url.URL
}

func NewInMemoryStorage() InMemoryStorage {
	return InMemoryStorage{
		shid2url: make(map[string]*url.URL),
	}
}

func (s InMemoryStorage) Put(shid string, url *url.URL) bool {
	s.shid2url[shid] = url
	return true
}

func (s InMemoryStorage) Get(shid string) (url *url.URL, isOk bool) {
	url, isOk = s.shid2url[shid]
	return
}
