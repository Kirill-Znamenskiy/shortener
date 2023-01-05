package storage

import (
	"database/sql"
	"errors"
	"github.com/Kirill-Znamenskiy/Shortener/internal/blogic/btypes"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/url"
)

type DBStorage struct {
	dsn    string
	dbconn *sql.DB
	pstmts map[string]*sql.Stmt
}

func NewDBStorage(dsn string) (ret *DBStorage, err error) {
	ret = &DBStorage{
		dsn:    dsn,
		pstmts: make(map[string]*sql.Stmt),
	}

	ret.dbconn, err = sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = ret.init()
	if err != nil {
		return nil, err
	}

	return
}

func (s *DBStorage) init() (err error) {
	_, err = s.dbconn.Exec(`
		CREATE TABLE IF NOT EXISTS records(
			key VARCHAR(32) PRIMARY KEY,
            original_url TEXT NOT NULL,
			usr UUID NOT NULL
		);
	`)
	if err != nil {
		return
	}

	_, err = s.dbconn.Exec(`
		CREATE TABLE IF NOT EXISTS secret_keys(
			name VARCHAR(32) PRIMARY KEY,
			value BYTEA
		);
	`)
	if err != nil {
		return
	}

	return
}

func (s *DBStorage) TruncateAllRecords() (err error) {
	_, err = s.dbconn.Exec(`TRUNCATE records`)
	if err != nil {
		return
	}

	return
}

func (s *DBStorage) PutSecretKey(secretKey []byte) (err error) {
	_, err = s.dbconn.Exec(`
		INSERT INTO secret_keys
			(name, value)
		VALUES
			($1, $2)
		ON CONFLICT (name) DO
			UPDATE SET value = EXCLUDED.value
	`, "main", secretKey)
	return
}

func (s *DBStorage) GetSecretKey() (secretKey []byte, err error) {
	row := s.dbconn.QueryRow(`
		SELECT value FROM secret_keys WHERE (name = $1)
	`, "main")
	err = row.Err()
	if err != nil {
		return
	}
	err = row.Scan(&secretKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make([]byte, 0), nil
		}
		return nil, err
	}
	return
}

func (s *DBStorage) PutRecord(record *btypes.Record) (err error) {
	_, err = s.dbconn.Exec(`
		INSERT INTO records
		    (key, original_url, usr)
		VALUES
		    ($1, $2, $3)
	`, record.Key, record.OriginalURL, record.User)
	return
}

func (s *DBStorage) GetRecord(key string) (rec *btypes.Record, err error) {
	key2record, err := s.extractKey2RecordFromRows(s.dbconn.Query(`
		SELECT key, original_url, usr FROM records WHERE (key = $1)
	`, key))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return
	}

	if len(key2record) == 0 {
		return nil, nil
	}

	for _, rec := range key2record {
		return rec, nil
	}

	return nil, errors.New("unexpected behavior")
}

func (s *DBStorage) extractKey2RecordFromRows(rows *sql.Rows, rowsErr error) (key2record map[string]*btypes.Record, err error) {
	if rowsErr != nil {
		err = rowsErr
		return
	}
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		return
	}

	key2record = make(map[string]*btypes.Record)
	var urlStr, usrStr string
	var tmpUUID uuid.UUID
	for rows.Next() {
		r := new(btypes.Record)
		err = rows.Scan(&r.Key, &urlStr, &usrStr)
		if err != nil {
			return nil, err
		}
		r.OriginalURL, err = url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		tmpUUID, err = uuid.Parse(usrStr)
		if err != nil {
			return nil, err
		}
		r.User = &tmpUUID

		key2record[r.Key] = r
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return
}

func (s *DBStorage) GetAllUserRecords(user btypes.User) (userKey2Record map[string]*btypes.Record, err error) {
	return s.extractKey2RecordFromRows(s.dbconn.Query(`
		SELECT key, original_url, usr FROM records WHERE (usr = $1)
	`, user))
}

func (s *DBStorage) Ping() error {
	return s.dbconn.Ping()
}
