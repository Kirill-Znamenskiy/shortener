package storage

import "net/url"

type Storage interface {
	Put(shid string, url *url.URL) bool
	Get(shid string) (*url.URL, bool)
}
