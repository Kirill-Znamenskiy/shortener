package storage

import (
	"github.com/Kirill-Znamenskiy/Shortener/internal/blogic/types"
)

type Storage interface {
	PutSecretKey(secretKey []byte) error
	GetSecretKey() []byte
	PutRecord(r *types.Record) error
	GetRecord(key string) (r *types.Record)
	GetAllUserRecords(user types.User) (userKey2Record map[string]*types.Record)
	Ping() error
}
