package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/viacheslaev/url-shortener/internal/config"
	"github.com/viacheslaev/url-shortener/internal/feature/link"
	"github.com/viacheslaev/url-shortener/internal/server"
	"github.com/viacheslaev/url-shortener/internal/server/middleware"
	"github.com/viacheslaev/url-shortener/internal/storage/postgres"
)

func main() {
	cfg := config.Load()

	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect db: %v", err)
	}
	defer disconnectDB(db)

	repo := postgres.NewLinkRepository(db)
	linkCfg := createLinkConfig(cfg)
	service := link.NewURLService(repo, linkCfg)
	handler := link.NewURLHandler(cfg, service)

	router := middleware.Logging(server.NewRouter(cfg, handler))

	log.Printf("listening on %s\n", cfg.HTTPAddr)
	log.Fatal(http.ListenAndServe(cfg.HTTPAddr, router))
}

func connectDB(cfg *config.Config) (*sql.DB, error) {
	return postgres.New(postgres.DBConfig{
		DSN: cfg.DSN,
	})
}

func disconnectDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("db close error: %v", err)
	}
}

func createLinkConfig(cfg *config.Config) *link.Config {
	return &link.Config{
		ShortLinkTTL: time.Duration(cfg.LinkTTLHours) * time.Second,
	}
}
