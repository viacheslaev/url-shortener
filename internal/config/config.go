package config

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr string
	BaseURL  string
	LogLevel string
}

func Load() *Config {
	if os.Getenv("APP_ENV") == "dev" {
		readEnvFile()
	}

	cfg := &Config{
		HTTPAddr: getEnv("HTTP_ADDR", ":8080"),
		BaseURL:  getEnv("BASE_URL", "http://localhost:8080"),
		LogLevel: getEnv("LOG_LEVEL", "INFO"),
	}

	validate(cfg)
	log.Printf("Config loaded APP_ENV=%s, LOG_LEVEL=%s\n", os.Getenv("APP_ENV"), cfg.LogLevel)
	return cfg
}

// readEnvFile loads .env variables from file when APP_ENV=dev
func readEnvFile() {
	log.Println("Reading properties from .env file")
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, loading default cfg")
		return
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
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

	// LogLevel
	switch strings.ToLower(cfg.LogLevel) {
	case "debug", "info", "warn", "error":
		cfg.LogLevel = strings.ToLower(cfg.LogLevel)
	default:
		log.Fatalf("LOG_LEVEL must be one of: debug, info, warn, error (got %q)", cfg.LogLevel)
	}
}
