package storage

import (
	"database/sql"
	"github.com/Kirill-Znamenskiy/Shortener/internal/blogic/btypes"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	dsn    string
	dbconn *sql.DB
}

func NewDBStorage(dsn string) (ret *DBStorage, err error) {
	ret = &DBStorage{
		dsn: dsn,
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

func (s *DBStorage) PutSecretKey(secretKey []byte) (err error) {
	_, err = s.dbconn.Exec(`
		INSERT INTO secret_keys
		    (name, value) 
		VALUES
		    (?, ?)
		ON CONFLICT (name) DO 
			UPDATE SET value = EXCLUDED.value
	`, "main", secretKey)
	return
}

func (s *DBStorage) GetSecretKey() (secretKey []byte, err error) {
	err = s.dbconn.QueryRow(
		`SELECT * FROM secret_keys WHERE (name = ?)`,
		"main",
	).Scan(&secretKey)
	if err != nil {
		return nil, err
	}
	return
}

func (s *DBStorage) PutRecord(record *btypes.Record) (err error) {
	_, err = s.dbconn.Exec(`
		INSERT INTO records 
		    (key, original_url, usr) 
		VALUES
		    (?, ?, ?)
	`, record.Key, record.OriginalURL, record.User)
	return
}

func (s *DBStorage) GetRecord(key string) (r *btypes.Record, err error) {
	r = new(btypes.Record)
	err = s.dbconn.QueryRow(
		`SELECT key, original_url, usr FROM records WHERE (key = ?)`,
		key,
	).Scan(&r.Key, &r.OriginalURL, &r.User)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s *DBStorage) GetAllUserRecords(user btypes.User) (userKey2Record map[string]*btypes.Record, err error) {
	userKey2Record = make(map[string]*btypes.Record)

	rows, err := s.dbconn.Query(
		`SELECT key, original_url, usr FROM records WHERE (usr = ?)`,
		user,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		r := new(btypes.Record)
		err = rows.Scan(&r.Key, &r.OriginalURL, &r.User)
		if err != nil {
			return nil, err
		}
		userKey2Record[r.Key] = r
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return
}

func (s *DBStorage) Ping() error {
	return s.dbconn.Ping()
}
