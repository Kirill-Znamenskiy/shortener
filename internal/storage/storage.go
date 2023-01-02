package storage

import "net/url"

type Storage interface {
	Put(key string, url *url.URL) error
	Get(key string) (url *url.URL, isOk bool)
}
