package config

import (
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr                         string
	BaseURL                          string
	LogLevel                         string
	DSN                              string
	LinkTTLHours                     int
	ExpiredLinksCleanupIntervalHours int
}

func Load() *Config {
	if os.Getenv("APP_ENV") == "dev" {
		readEnvFile()
	}

	cfg := &Config{
		HTTPAddr:                         getEnv("HTTP_ADDR"),
		BaseURL:                          getEnv("BASE_URL"),
		LogLevel:                         getEnv("LOG_LEVEL"),
		DSN:                              getEnv("DSN"),
		LinkTTLHours:                     getEnvInt("LINK_TTL_HOURS"),
		ExpiredLinksCleanupIntervalHours: getEnvInt("EXPIRED_LINKS_CLEANUP_INTERVAL_HOURS"),
	}

	validate(cfg)

	log.Printf("config loaded APP_ENV=%s, LOG_LEVEL=%s\n", os.Getenv("APP_ENV"), cfg.LogLevel)
	return cfg
}

// readEnvFile loads .env variables from file when APP_ENV=dev
func readEnvFile() {
	log.Println("reading properties from .env file")
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error reading .env file %q", err)
	}
}

func getEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		log.Fatalf("env %s is required", key)
	}
	return v
}

func validate(cfg *Config) {
	// HTTPAddr
	if strings.TrimSpace(cfg.HTTPAddr) == "" {
		log.Fatal("HTTP_ADDR is required")
	}

	// BaseURL
	if strings.TrimSpace(cfg.BaseURL) == "" {
		log.Fatal("BASE_URL is required")
	}

	u, err := url.Parse(cfg.BaseURL)
	if err != nil || !u.IsAbs() || (u.Scheme != "http" && u.Scheme != "https") {
		log.Fatalf("BASE_URL must be a valid http/https URL: %q", cfg.BaseURL)
	}

	// normalize: remove trailing slash
	cfg.BaseURL = strings.TrimRight(cfg.BaseURL, "/")

	// DSN
	if strings.TrimSpace(cfg.DSN) == "" {
		log.Fatal("DSN is required")
	}

	// LINK
	if cfg.LinkTTLHours <= 0 {
		log.Fatalf("LINK_TTL_HOURS must be > 0 (got %d)", cfg.LinkTTLHours)
	}
	if cfg.ExpiredLinksCleanupIntervalHours <= 0 {
		log.Fatalf("EXPIRED_LINKS_CLEANUP_INTERVAL_HOURS must be > 0 (got %d)", cfg.LinkTTLHours)
	}

	// LogLevel
	switch strings.ToLower(cfg.LogLevel) {
	case "debug", "info", "warn", "error":
		cfg.LogLevel = strings.ToLower(cfg.LogLevel)
	default:
		log.Fatalf("LOG_LEVEL must be one of: debug, info, warn, error (got %q)", cfg.LogLevel)
	}
}

func getEnvInt(key string) int {
	v := getEnv(key)
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("env %s must be an integer (got %q)", key, v)
	}
	return i
}
