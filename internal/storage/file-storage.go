package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/url"
	"os"
)

type FileStorage struct {
	*InMemoryStorage
	filePath string
}

func NewFileStorage(filePath string) (ret *FileStorage) {
	ret = &FileStorage{
		filePath:        filePath,
		InMemoryStorage: NewInMemoryStorage(),
	}
	err := ret.LoadDataFromFile()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (s *FileStorage) Put(key string, url *url.URL) (retKey string, err error) {
	retKey, err = s.InMemoryStorage.Put(key, url)
	if err == nil {
		return
	}
	err = s.SaveDataToFile()
	return
}

func (s *FileStorage) LoadDataFromFile() (err error) {
	s.InMemoryStorage = NewInMemoryStorage()

	file, err := os.Open(s.filePath)
	if err == os.ErrNotExist {
		return nil
	}
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
		err = json.Unmarshal(toLoadBytes, &s.InMemoryStorage.key2url)
		if err != nil {
			return
		}
	}

	return
}

func (s *FileStorage) SaveDataToFile() (err error) {
	toSaveData := s.InMemoryStorage.key2url
	if len(toSaveData) == 0 {
		return
	}

	file, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	toSaveBytes, err := json.Marshal(&toSaveData)
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
