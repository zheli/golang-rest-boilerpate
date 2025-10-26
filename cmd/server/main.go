package main

import (
	"fmt"
	"log"

	"github.com/example/golang-rest-boilerplate/internal/config"
	"github.com/example/golang-rest-boilerplate/internal/db"
	"github.com/example/golang-rest-boilerplate/internal/http/handlers"
	"github.com/example/golang-rest-boilerplate/internal/http/router"
	"github.com/example/golang-rest-boilerplate/internal/repository"
	"github.com/example/golang-rest-boilerplate/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	userRepo := repository.NewUserRepository(database)
	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo)

	var googleService *service.GoogleOAuthService
	if cfg.GoogleClientID != "" && cfg.GoogleClientSecret != "" {
		googleService = service.NewGoogleOAuthService(cfg)
	}

	authHandler := handlers.NewAuthHandler(authService, googleService)
	userHandler := handlers.NewUserHandler(userService)
	healthHandler := handlers.NewHealthHandler()

	r := router.SetupRouter(authHandler, userHandler, healthHandler, authService, cfg)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
