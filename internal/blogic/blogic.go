package blogic

import (
	"github.com/Kirill-Znamenskiy/Shortener/internal/storage"
	"github.com/teris-io/shortid"
	"net/url"
)

// SaveNewURL save new url in storage.
func SaveNewURL(stg storage.Storage, urlStr string) (key string, err error) {
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		return
	}

	key, err = GenerateNewKey(stg)
	if err != nil {
		return
	}

	err = stg.Put(key, urlObj)

	return
}

// GetSavedURL extract early saved url from storage by key.
func GetSavedURL(stg storage.Storage, key string) (u *url.URL, isOk bool) {
	return stg.Get(key)
}

// GenerateNewKey generate new key, that don't exist in storage.
func GenerateNewKey(stg storage.Storage) (ret string, err error) {

	for {
		ret, err = shortid.Generate()
		if err != nil {
			return
		}
		if _, isOk := stg.Get(ret); isOk {
			continue
		}
		break
	}

	return
}
