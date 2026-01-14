package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
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
	linkService := link.NewLinkService(analyticsService, linkRepo, cfg)
	accountService := account.NewAccountService(accountRepo)

	// AUTH (register/login + JWT)
	tokenIssuer := auth.NewTokenIssuer(cfg)
	authService := auth.NewAuthService(accountRepo, tokenIssuer)
	authMiddleware := middleware.NewAuthMiddleware(accountRepo, cfg)

	// JOB
	expiredLinksCleanupWorker := link.NewExpiredLinksCleanupWorker(linkRepo, time.Duration(cfg.ExpiredLinksCleanupIntervalHours)*time.Hour)
	expiredLinksCleanupWorker.Start()
	clickEventWorker := analytics.NewClickEventWorker(analyticsService)
	clickEventWorker.Start()

	// HANDLER
	linkHandler := link.NewLinkHandler(cfg, linkService)
	accRegisterHandler := account.NewAccountRegisterHandler(accountService)
	authHandler := auth.NewAuthHandler(authService)
	analyticsHandler := analytics.NewAnalyticsHandler(linkRepo, analyticsService)

	// ROUTER
	router := middleware.Logging(server.NewRouter(linkHandler, accRegisterHandler, authHandler, analyticsHandler, authMiddleware))

	// SERVER
	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}

	// Start server
	go func() {
		log.Printf("starting server on port%s\n", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server failed: %v", err)
		}
	}()

	// SHUTDOWN CONTEXT
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// GRACEFULLY SHUTDOWN
	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	clickEventWorker.Stop()
	expiredLinksCleanupWorker.Stop()

	log.Println("server stopped")
}
