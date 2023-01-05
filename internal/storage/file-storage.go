package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/Kirill-Znamenskiy/Shortener/internal/blogic/types"
	"io"
	"log"
	"os"
)

type FileStorage struct {
	*InMemoryStorage
	filePath string
}

func NewFileStorage(filePath string) (ret *FileStorage) {
	ret = &FileStorage{
		filePath: filePath,
	}
	err := ret.LoadDataFromFile()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (s *FileStorage) PutSecretKey(secretKey []byte) (err error) {
	err = s.InMemoryStorage.PutSecretKey(secretKey)
	if err != nil {
		return
	}
	err = s.SaveDataToFile()
	return
}
func (s *FileStorage) PutRecord(r *types.Record) (err error) {
	err = s.InMemoryStorage.PutRecord(r)
	if err != nil {
		return
	}
	err = s.SaveDataToFile()
	return
}

func (s *FileStorage) LoadDataFromFile() (err error) {
	s.InMemoryStorage = NewInMemoryStorage()

	file, err := os.Open(s.filePath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	toLoadBytes, err := io.ReadAll(reader)
	if err != nil {
		return
	}

	if len(toLoadBytes) > 0 {
		err = json.Unmarshal(toLoadBytes, s.InMemoryStorage)
		if err != nil {
			return
		}
	}

	return
}

func (s *FileStorage) SaveDataToFile() (err error) {

	file, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	if s.InMemoryStorage.IsEmpty() {
		return
	}

	toSaveBytes, err := json.Marshal(s.InMemoryStorage)
	if err != nil {
		return
	}

	writer := bufio.NewWriter(file)
	_, err = writer.Write(toSaveBytes)
	if err != nil {
		return
	}

	err = writer.Flush()
	if err != nil {
		return
	}

	return
}
