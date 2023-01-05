package blogic

import (
	"github.com/Kirill-Znamenskiy/Shortener/internal/blogic/types"
	"github.com/Kirill-Znamenskiy/Shortener/internal/storage"
	"github.com/google/uuid"
	"github.com/teris-io/shortid"
	"net/url"
)

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

// BuildShortURL make record short url, by record key.
func (sh *Shortener) BuildShortURL(record *types.Record) (recordShortURL string) {
	return sh.baseURL + "/" + record.Key
}

// SaveNewURL save new url in storage.
func (sh *Shortener) SaveNewURL(user *uuid.UUID, urlStr string) (ret *types.Record, err error) {
	record := new(types.Record)
	record.User = user

	record.OriginalURL, err = url.Parse(urlStr)
	if err != nil {
		return
	}

	record.Key, err = sh.GenerateNewRecordKey()
	if err != nil {
		return
	}

	err = sh.stg.PutRecord(record)
	if err != nil {
		return
	}

	ret = record
	return
}

// GetSavedURL extract early saved url from storage by key.
func (sh *Shortener) GetSavedURL(recordKey string) (u *url.URL, isOk bool) {
	record := sh.stg.GetRecord(recordKey)
	if record == nil {
		return
	}
	return record.OriginalURL, true
}

func (sh *Shortener) GetAllUserRecords(user types.User) (ret map[string]*types.Record) {
	return sh.stg.GetAllUserRecords(user)
}

// GenerateNewRecordKey generate new record key, that don't exist in storage.
func (sh *Shortener) GenerateNewRecordKey() (ret string, err error) {

	for {
		ret, err = shortid.Generate()
		if err != nil {
			return
		}
		if r := sh.stg.GetRecord(ret); r == nil {
			break
		}
	}

	return
}

// GenerateNewUser generate new user UUID, that don't exist in storage.
func (sh *Shortener) GenerateNewUser() (ret types.User, err error) {

	var tmp uuid.UUID

	for {
		tmp, err = uuid.NewRandom()
		if err != nil {
			return
		}
		ret = &tmp
		if m := sh.stg.GetAllUserRecords(ret); len(m) == 0 {
			break
		}
	}

	return
}
