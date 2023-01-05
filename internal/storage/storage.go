package storage

import (
	"github.com/Kirill-Znamenskiy/Shortener/internal/blogic/btypes"
)

type Storage interface {
	PutSecretKey(secretKey []byte) error
	GetSecretKey() (secretKey []byte, err error)
	PutRecord(r *btypes.Record) error
	PutRecords(r []*btypes.Record) error
	GetRecord(key string) (r *btypes.Record, err error)
	GetAllUserRecords(user btypes.User) (userKey2Record map[string]*btypes.Record, err error)
	Ping() error
}
