package storage

import "net/url"

type Storage interface {
	Put(key string, url *url.URL) (string, error)
	Get(key string) (*url.URL, bool)
	GenerateNewKey() (string, error)
}
