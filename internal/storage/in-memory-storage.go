package storage

import (
	"fmt"
	"github.com/google/uuid"
	"net/url"
)

type InMemoryStorage struct {
	secretKey              []byte
	userUUID2RecordKey2URL map[uuid.UUID]map[string]*url.URL
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		userUUID2RecordKey2URL: make(map[uuid.UUID]map[string]*url.URL),
	}
}

func (s *InMemoryStorage) PutSecretKey(secretKey []byte) (err error) {
	s.secretKey = secretKey
	return nil
}

func (s *InMemoryStorage) GetSecretKey() []byte {
	return s.secretKey
}

func (s *InMemoryStorage) Put(userUUID *uuid.UUID, recordKey string, urlParam *url.URL) (err error) {
	recordKey2URL, isOk := s.userUUID2RecordKey2URL[*userUUID]
	if isOk {
		if _, isOk := recordKey2URL[recordKey]; isOk {
			return fmt.Errorf("recordKey %q already exists", recordKey)
		}
	} else {
		s.userUUID2RecordKey2URL[*userUUID] = make(map[string]*url.URL)
	}

	s.userUUID2RecordKey2URL[*userUUID][recordKey] = urlParam
	return
}

func (s *InMemoryStorage) Get(userUUID *uuid.UUID, recordKey string) (url *url.URL, isOk bool) {
	recordKey2URL, isOk := s.GetAllUserURLs(userUUID)
	if !isOk {
		return nil, false
	}

	url, isOk = recordKey2URL[recordKey]
	return
}

func (s *InMemoryStorage) GetAllUserURLs(userUUID *uuid.UUID) (userRecordKey2URL map[string]*url.URL, isOk bool) {
	userRecordKey2URL, isOk = s.userUUID2RecordKey2URL[*userUUID]
	return
}

func (s *InMemoryStorage) setSrcMap(mp map[uuid.UUID]map[string]*url.URL) {
	s.userUUID2RecordKey2URL = mp
}
func (s *InMemoryStorage) getSrcMap() map[uuid.UUID]map[string]*url.URL {
	return s.userUUID2RecordKey2URL
}
