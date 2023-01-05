package storage

import (
	"fmt"
	"github.com/Kirill-Znamenskiy/Shortener/internal/blogic/types"
)

type InMemoryStorage struct {
	SecretKey  []byte
	Key2Record map[string]*types.Record
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		Key2Record: make(map[string]*types.Record),
	}
}

func (s *InMemoryStorage) PutSecretKey(secretKey []byte) (err error) {
	s.SecretKey = secretKey
	return nil
}

func (s *InMemoryStorage) GetSecretKey() []byte {
	return s.SecretKey
}

func (s *InMemoryStorage) PutRecord(r *types.Record) (err error) {
	if _, isAlreadyExists := s.Key2Record[r.Key]; isAlreadyExists {
		return fmt.Errorf("record with key %q already exists", r.Key)
	}

	s.Key2Record[r.Key] = r
	return
}

func (s *InMemoryStorage) GetRecord(key string) (r *types.Record) {
	r, isOk := s.Key2Record[key]
	if !isOk {
		return nil
	}

	return
}

func (s *InMemoryStorage) GetAllUserRecords(user types.User) (userKey2Record map[string]*types.Record) {
	userKey2Record = make(map[string]*types.Record)
	for key, record := range s.Key2Record {
		if *record.User == *user {
			userKey2Record[key] = record
		}
	}
	return
}

func (s *InMemoryStorage) IsEmpty() bool {
	return len(s.SecretKey) == 0 && len(s.Key2Record) == 0
}

func (s *InMemoryStorage) Ping() error {
	return nil
}
