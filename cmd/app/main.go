package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/viacheslaev/url-shortener/internal/config"
	"github.com/viacheslaev/url-shortener/internal/feature/account"
	"github.com/viacheslaev/url-shortener/internal/feature/analytics"
	"github.com/viacheslaev/url-shortener/internal/feature/auth"
	"github.com/viacheslaev/url-shortener/internal/feature/link"
	"github.com/viacheslaev/url-shortener/internal/server"
	"github.com/viacheslaev/url-shortener/internal/server/middleware"
	"github.com/viacheslaev/url-shortener/internal/storage/postgres"
)

func main() {
	// CONFIG
	cfg := config.Load()
	linkCfg := createLinkConfig(cfg)

	// DB
	db, err := postgres.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer postgres.DisconnectDB(db)

	// REPOSITORY
	linkRepo := postgres.NewLinkRepository(db)
	accountRepo := postgres.NewAccountRepository(db)
	analyticsRepository := postgres.NewAnalyticsRepository(db)

	// SERVICE
	analyticsService := analytics.NewAnalyticsService(analyticsRepository)
	urlService := link.NewURLService(analyticsService, linkRepo, linkCfg)
	accountService := account.NewAccountService(accountRepo)

	// AUTH (register/login + JWT)
	tokenIssuer := auth.NewTokenIssuer(cfg)
	authService := auth.NewAuthService(accountRepo, tokenIssuer)
	authMiddleware := middleware.NewAuthMiddleware(accountRepo, cfg)

	// HANDLER
	urlHandler := link.NewURLHandler(cfg, urlService)
	accRegisterHandler := account.NewAccountRegisterHandler(accountService)
	authHandler := auth.NewAuthHandler(authService)

	// ROUTER
	router := middleware.Logging(server.NewRouter(urlHandler, accRegisterHandler, authHandler, authMiddleware))

	// SERVER
	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}

	// SHUTDOWN CONTEXT
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// JOB
	cleanupJobStop, cleanupJobWait := startExpiredLinksCleanupJob(linkRepo, time.Duration(cfg.ExpiredLinksCleanupIntervalHours)*time.Hour)

	// Start server
	go func() {
		log.Printf("starting server on port%s\n", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server failed: %v", err)
		}
	}()

	// GRACEFULLY SHUTDOWN
	<-ctx.Done()
	log.Println("shutdown signal received")

	cleanupJobStop()
	cleanupJobWait()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	} else {
		log.Println("server stopped")
	}
}

func createLinkConfig(cfg *config.Config) *link.Config {
	return &link.Config{
		ShortLinkTTL: time.Duration(cfg.LinkTTLHours) * time.Hour,
	}
}

// startExpiredLinksCleanupJob starts background job that periodically deletes expired links.
// It returns two functions:
//   - stop(): signals the worker to stop
//   - wait(): blocks until the worker finishes cleanup job
func startExpiredLinksCleanupJob(repo link.LinkRepository, interval time.Duration) (stop func(), wait func()) {
	ticker := time.NewTicker(interval)
	done := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		log.Printf("started expired links cleanup job: interval=%s", interval)

		for {
			select {
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				deleted, err := repo.DeleteExpiredLinks(ctx)
				cancel()

				if err != nil {
					log.Printf("cleanup expired links failed: %v", err)
				} else if deleted > 0 {
					log.Printf("cleanup expired links: deleted=%d", deleted)
				}

			case <-done:
				ticker.Stop()
				log.Printf("stopped expired links cleanup job")
				return
			}
		}
	}()

	return func() { close(done) }, wg.Wait
}
