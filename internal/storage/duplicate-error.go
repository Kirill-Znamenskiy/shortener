package storage

import (
	"errors"
	"github.com/Kirill-Znamenskiy/Shortener/internal/blogic/btypes"
)

var ErrDupOriginalURL = errors.New("duplicated original url")

type DuplicateError struct {
	Err           error
	NewRecord     *btypes.Record
	AlreadyRecord *btypes.Record
}

func NewDuplicateError(newRecord *btypes.Record, alreadyRecord *btypes.Record) error {
	return &DuplicateError{
		Err:           ErrDupOriginalURL,
		NewRecord:     newRecord,
		AlreadyRecord: alreadyRecord,
	}
}

func (de *DuplicateError) Error() string {
	return de.Err.Error()
}

func (de *DuplicateError) Unwrap() error {
	return de.Err
}
