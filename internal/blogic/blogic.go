package blogic

import (
	"github.com/Kirill-Znamenskiy/Shortener/internal/blogic/btypes"
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
func (sh *Shortener) BuildShortURL(record *btypes.Record) (recordShortURL string) {
	return sh.baseURL + "/" + record.Key
}

// SaveNewURL save new url in storage.
func (sh *Shortener) SaveNewURL(user *uuid.UUID, urlStr string) (ret *btypes.Record, err error) {
	tmp, err := sh.BatchSaveNewURL(user, []string{urlStr})
	if err != nil {
		return
	}
	ret = tmp[0]
	return
}

// BatchSaveNewURL save batch or new urls in storage.
func (sh *Shortener) BatchSaveNewURL(user *uuid.UUID, urls []string) (records []*btypes.Record, err error) {
	records = make([]*btypes.Record, len(urls))
	for ind, urlStr := range urls {

		record := new(btypes.Record)
		record.User = user

		record.OriginalURL, err = url.Parse(urlStr)
		if err != nil {
			return nil, err
		}

		record.Key, err = sh.GenerateNewRecordKey()
		if err != nil {
			return nil, err
		}

		records[ind] = record
	}

	err = sh.stg.PutRecords(records)
	if err != nil {
		return nil, err
	}

	return
}

// GetSavedURL extract early saved url from storage by key.
func (sh *Shortener) GetSavedURL(recordKey string) (u *url.URL, isOk bool) {
	record, _ := sh.stg.GetRecord(recordKey)
	if record == nil {
		return
	}
	return record.OriginalURL, true
}

func (sh *Shortener) GetAllUserRecords(user btypes.User) (ret map[string]*btypes.Record, err error) {
	return sh.stg.GetAllUserRecords(user)
}

// GenerateNewRecordKey generate new record key, that don't exist in storage.
func (sh *Shortener) GenerateNewRecordKey() (ret string, err error) {

	for {
		ret, err = shortid.Generate()
		if err != nil {
			return
		}
		if r, _ := sh.stg.GetRecord(ret); r == nil {
			break
		}
	}

	return
}

// GenerateNewUser generate new user UUID, that don't exist in storage.
func (sh *Shortener) GenerateNewUser() (ret btypes.User, err error) {

	var tmp uuid.UUID

	for {
		tmp, err = uuid.NewRandom()
		if err != nil {
			return
		}
		ret = &tmp
		if m, _ := sh.stg.GetAllUserRecords(ret); len(m) == 0 {
			break
		}
	}

	return
}
