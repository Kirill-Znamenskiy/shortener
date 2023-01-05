package btypes

import (
	"github.com/google/uuid"
	"net/url"
)

type User *uuid.UUID

type Record struct {
	Key         string
	OriginalURL *url.URL
	User        User
}
