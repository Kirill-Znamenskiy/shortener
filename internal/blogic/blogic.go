package blogic

import (
	"github.com/Kirill-Znamenskiy/Shortener/internal/storage"
	"github.com/google/uuid"
	"github.com/teris-io/shortid"
	"net/url"
)

type User *uuid.UUID

// Shortener is main shortener struct.
type Shortener struct {
	baseURL string
	stg     storage.Storage
}

// NewShortener returns a new Shortener that uses specified storage.
func NewShortener(baseURL string, stg storage.Storage) *Shortener {
	return &Shortener{
		baseURL: baseURL,
		stg:     stg,
	}
}

// MakeShortURL make record short url, by record key.
func (sh *Shortener) MakeShortURL(recordKey string) (recordShortURL string) {
	return sh.baseURL + "/" + recordKey
}

// SaveNewURL save new url in storage.
func (sh *Shortener) SaveNewURL(userUUID *uuid.UUID, urlStr string) (shortURL string, err error) {
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		return
	}

	recordKey, err := sh.GenerateNewRecordKey(userUUID)
	if err != nil {
		return
	}

	err = sh.stg.Put(userUUID, recordKey, urlObj)
	if err != nil {
		return
	}

	return sh.MakeShortURL(recordKey), nil
}

// GetSavedURL extract early saved url from storage by key.
func (sh *Shortener) GetSavedURL(userUUID *uuid.UUID, recordKey string) (u *url.URL, isOk bool) {
	return sh.stg.Get(userUUID, recordKey)
}

func (sh *Shortener) GetAllUserURLs(userUUID *uuid.UUID) (userRecordKey2URL map[string]*url.URL) {
	userRecordKey2URL, _ = sh.stg.GetAllUserURLs(userUUID)
	return userRecordKey2URL
}

// GenerateNewRecordKey generate new record key, that don't exist in storage.
func (sh *Shortener) GenerateNewRecordKey(userUUID *uuid.UUID) (ret string, err error) {

	for {
		ret, err = shortid.Generate()
		if err != nil {
			return
		}
		if _, isOk := sh.stg.Get(userUUID, ret); isOk {
			continue
		}
		break
	}

	return
}

// GenerateNewUserUUID generate new user UUID, that don't exist in storage.
func (sh *Shortener) GenerateNewUserUUID() (ret *uuid.UUID, err error) {

	var tmp uuid.UUID

	for {
		tmp, err = uuid.NewRandom()
		if err != nil {
			return
		}
		ret = &tmp
		if _, isOk := sh.stg.GetAllUserURLs(ret); isOk {
			continue
		}
		break
	}

	return
}
