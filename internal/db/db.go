package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/example/golang-rest-boilerplate/internal/config"
	"github.com/example/golang-rest-boilerplate/internal/models"
)

// New initializes a new GORM database instance using configuration.
func New(cfg *config.Config) (*gorm.DB, error) {
	database, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := database.AutoMigrate(&models.User{}); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return database, nil
}
