package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string
	DatabaseURL            string
	JWTSecret              string
	JWTExpiration          time.Duration
	RefreshTokenExpiration time.Duration
	OpenLibraryURL         string
	GoogleBooksURL         string
	GoogleBooksAPIKey      string
	AllowedOrigin          string
	RequestTimeoutSec      int
	FrontendURL            string
	SMTPHost               string
	SMTPPort               string
	SMTPUsername           string
	SMTPPassword           string
	SMTPFrom               string
	RefreshCookieSecure    bool
}

func Load() Config {
	_ = godotenv.Load()
	exp := getOrDefault("JWT_EXPIRATION", "15m")
	duration, err := time.ParseDuration(exp)
	if err != nil || duration <= 0 {
		duration = 15 * time.Minute
	}
	refreshExp := getOrDefault("REFRESH_TOKEN_EXPIRATION", "720h")
	refreshDuration, err := time.ParseDuration(refreshExp)
	if err != nil || refreshDuration <= 0 {
		refreshDuration = 30 * 24 * time.Hour
	}

	timeoutStr := getOrDefault("REQUEST_TIMEOUT_SECONDS", "8")
	timeoutSec, err := strconv.Atoi(timeoutStr)
	if err != nil {
		timeoutSec = 8
	}

	return Config{
		Port:                   getOrDefault("PORT", "3000"),
		DatabaseURL:            os.Getenv("DATABASE_URL"),
		JWTSecret:              getOrDefault("JWT_SECRET", "segredo-dev-mude-isso"),
		JWTExpiration:          duration,
		RefreshTokenExpiration: refreshDuration,
		OpenLibraryURL:         getOrDefault("OPEN_LIBRARY_BASE_URL", "https://openlibrary.org"),
		GoogleBooksURL:         getOrDefault("GOOGLE_BOOKS_BASE_URL", "https://www.googleapis.com/books/v1"),
		GoogleBooksAPIKey:      os.Getenv("GOOGLE_BOOKS_API_KEY"),
		AllowedOrigin:          getOrDefault("CORS_ALLOWED_ORIGIN", "http://localhost:5173"),
		RequestTimeoutSec:      timeoutSec,
		FrontendURL:            getOrDefault("FRONTEND_URL", "http://localhost:5173"),
		SMTPHost:               os.Getenv("SMTP_HOST"),
		SMTPPort:               getOrDefault("SMTP_PORT", "587"),
		SMTPUsername:           os.Getenv("SMTP_USERNAME"),
		SMTPPassword:           os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:               os.Getenv("SMTP_FROM"),
		RefreshCookieSecure:    getOrDefault("REFRESH_COOKIE_SECURE", "false") == "true",
	}
}

func getOrDefault(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
