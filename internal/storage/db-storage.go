package storage

import (
	"context"
	"github.com/Kirill-Znamenskiy/Shortener/internal/blogic/types"
	"github.com/jackc/pgx/v5"
	"log"
)

type DBStorage struct {
	dsn  string
	ctx  context.Context
	conn *pgx.Conn
}

func NewDBStorage(ctx context.Context, dsn string) (ret *DBStorage) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	return &DBStorage{
		dsn:  dsn,
		ctx:  ctx,
		conn: conn,
	}
}

func (s *DBStorage) PutSecretKey([]byte) (err error) {
	return nil
}

func (s *DBStorage) GetSecretKey() []byte {
	return nil
}

func (s *DBStorage) PutRecord(*types.Record) (err error) {
	return nil
}

func (s *DBStorage) GetRecord(string) (r *types.Record) {
	return nil
}

func (s *DBStorage) GetAllUserRecords(types.User) (userKey2Record map[string]*types.Record) {
	userKey2Record = make(map[string]*types.Record)
	return
}

func (s *DBStorage) Ping() error {
	return s.conn.Ping(s.ctx)
}
