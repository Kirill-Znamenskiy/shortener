package storage

import (
	"github.com/google/uuid"
	"net/url"
)

type Storage interface {
	PutSecretKey(secretKey []byte) error
	GetSecretKey() []byte
	Put(userUUID *uuid.UUID, recordKey string, url *url.URL) error
	Get(userUUID *uuid.UUID, recordKey string) (url *url.URL, isOk bool)
	GetAllUserURLs(userUUID *uuid.UUID) (userRecordKey2URL map[string]*url.URL, isOk bool)
}
