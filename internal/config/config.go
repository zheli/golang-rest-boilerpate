package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config holds configuration values for the application.
type Config struct {
	AppPort            string   `env:"APP_PORT" default:"8080"`
	DatabaseURL        string   `env:"DATABASE_URL" default:"postgres://postgres:postgres@db:5432/app?sslmode=disable"`
	JWTSecret          string   `env:"JWT_SECRET" default:"change-me"`
	JWTIssuer          string   `env:"JWT_ISSUER" default:"golang-rest-boilerplate"`
	TokenExpireMinutes int      `env:"TOKEN_EXPIRE_MINUTES" default:"60"`
	GoogleClientID     string   `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string   `env:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL  string   `env:"GOOGLE_REDIRECT_URL" default:"http://localhost:8080/api/v1/auth/google/callback"`
	AllowedOrigins     []string `env:"ALLOWED_ORIGINS" default:"*"`
}

// Load reads configuration from environment variables and .env files.
func Load() (*Config, error) {
	// Load .env if present, ignore error in production scenarios.
	if err := godotenv.Load(); err != nil {
		log.Printf("config: no .env file found: %v", err)
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
