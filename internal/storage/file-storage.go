package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
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
func (s *FileStorage) Put(userUUID *uuid.UUID, recordKey string, url *url.URL) (err error) {
	err = s.InMemoryStorage.Put(userUUID, recordKey, url)
	if err != nil {
		return
	}
	err = s.SaveDataToFile()
	return
}

type inFileData struct {
	SecretKey              []byte
	UserUUID2RecordKey2URL map[uuid.UUID]map[string]*url.URL
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
		var toLoadData inFileData
		err = json.Unmarshal(toLoadBytes, &toLoadData)
		if err != nil {
			return
		}
		err = s.InMemoryStorage.PutSecretKey(toLoadData.SecretKey)
		if err != nil {
			return
		}
		s.InMemoryStorage.setSrcMap(toLoadData.UserUUID2RecordKey2URL)
	}

	return
}

func (s *FileStorage) SaveDataToFile() (err error) {

	file, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	toSaveData := new(inFileData)
	toSaveData.SecretKey = s.InMemoryStorage.GetSecretKey()
	toSaveData.UserUUID2RecordKey2URL = s.InMemoryStorage.getSrcMap()
	if len(toSaveData.SecretKey) == 0 && len(toSaveData.UserUUID2RecordKey2URL) == 0 {
		return
	}

	toSaveBytes, err := json.Marshal(toSaveData)
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
