package main

import (
	"context"
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
	// CONFIGS
	cfg := config.Load()
	linkCfg := createLinkConfig(cfg)

	// DB
	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect db: %v", err)
	}
	defer disconnectDB(db)

	// REPOSITORY
	repo := postgres.NewLinkRepository(db)

	// SERVICES
	service := link.NewURLService(repo, linkCfg)
	handler := link.NewURLHandler(cfg, service)

	// ROUTER
	router := middleware.Logging(server.NewRouter(cfg, handler))

	// JOBS
	startExpiredLinksCleanupJob(repo, time.Duration(cfg.ExpiredLinksCleanupIntervalHours)*time.Hour)

	// SERVER
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
		ShortLinkTTL: time.Duration(cfg.LinkTTLHours) * time.Hour,
	}
}

func startExpiredLinksCleanupJob(repo link.Repository, interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		log.Printf("started expired links cleanup job: interval=%s", interval)
		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			deleted, err := repo.DeleteExpiredLinks(ctx)
			cancel()

			if err != nil {
				log.Printf("cleanup expired links failed: %v", err)
				continue
			}
			if deleted > 0 {
				log.Printf("cleanup expired links: deleted=%d", deleted)
			}
		}
	}()
}
