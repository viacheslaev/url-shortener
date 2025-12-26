package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	_ "github.com/lib/pq"
	"github.com/viacheslaev/url-shortener/internal/config"
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

func ConnectDB(cfg *config.Config) (*sql.DB, error) {
	db, err := New(DBConfig{
		DSN: cfg.DSN,
	})
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		DisconnectDB(db)
		return nil, fmt.Errorf("db ping failed: %w", err)
	}

	dbUrl, err := url.Parse(cfg.DSN)
	if err != nil {
		log.Println("postgres DSN parse failed")
		return db, nil
	}

	log.Printf("connected to Postgres on port:%v", dbUrl.Port())

	return db, nil
}

func DisconnectDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("db close error: %v", err)
	}

	log.Printf("postgres disconnected")
}
