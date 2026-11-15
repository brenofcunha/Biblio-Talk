package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	DatabaseURL         string
	JWTSecret           string
	JWTExpiration       time.Duration
	OpenLibraryURL      string
	GoogleBooksURL      string
	GoogleBooksAPIKey   string
	AllowedOrigin       string
	RequestTimeoutSec   int
}

func Load() Config {
	_ = godotenv.Load()
	exp := getOrDefault("JWT_EXPIRATION", "24h")
	duration, err := time.ParseDuration(exp)
	if err != nil {
		duration = 24 * time.Hour
	}

	timeoutStr := getOrDefault("REQUEST_TIMEOUT_SECONDS", "8")
	timeoutSec, err := strconv.Atoi(timeoutStr)
	if err != nil {
		timeoutSec = 8
	}

	return Config{
		Port:              getOrDefault("PORT", "3000"),
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		JWTSecret:         getOrDefault("JWT_SECRET", "segredo-dev-mude-isso"),
		JWTExpiration:     duration,
		OpenLibraryURL:    getOrDefault("OPEN_LIBRARY_BASE_URL", "https://openlibrary.org"),
		GoogleBooksURL:    getOrDefault("GOOGLE_BOOKS_BASE_URL", "https://www.googleapis.com/books/v1"),
		GoogleBooksAPIKey: os.Getenv("GOOGLE_BOOKS_API_KEY"),
		AllowedOrigin:     getOrDefault("CORS_ALLOWED_ORIGIN", "http://localhost:5173"),
		RequestTimeoutSec: timeoutSec,
	}
}

func getOrDefault(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
