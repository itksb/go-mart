package dbpgsql

import (
	"context"
	"github.com/itksb/go-mart/pkg/logger"
	"github.com/jmoiron/sqlx"
)

// PgStorage - abstract database service
type PgStorage struct {
	dsn string
	db  *sqlx.DB
	l   logger.Interface
}

// NewPostgres - postgres service constructor
func NewPostgres(dsn string, l logger.Interface) (*PgStorage, error) {
	db, err := sqlx.Connect("postgres", dsn)
	return &PgStorage{
		dsn: dsn,
		db:  db,
		l:   l,
	}, err
}

func (s *PgStorage) reconnect(ctx context.Context) error {
	var err error
	if s.db == nil {
		s.db, err = sqlx.ConnectContext(ctx, "postgres", s.dsn)
		if err != nil {
			return err
		}
	}
	return nil
}

// Ping - check whether connection to db is valid or not
func (s *PgStorage) Ping(ctx context.Context) bool {
	var err error
	err = s.reconnect(ctx)
	if err != nil {
		s.l.Errorf(err.Error())
		return false
	}
	err = s.db.PingContext(ctx)
	if err != nil {
		s.l.Errorf(err.Error())
		return false
	}
	return true
}

// Close - destructor
func (s *PgStorage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
