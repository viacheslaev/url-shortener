package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	DSN string
}

func New(cfg DBConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
