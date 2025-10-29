package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/example/golang-rest-boilerplate/internal/config"
	"github.com/example/golang-rest-boilerplate/internal/http/handlers"
	"github.com/example/golang-rest-boilerplate/internal/http/middleware"
	"github.com/example/golang-rest-boilerplate/internal/service"
)

// SetupRouter configures the gin router and routes.
func SetupRouter(authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, healthHandler *handlers.HealthHandler, authService *service.AuthService, cfg *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	if len(cfg.AllowedOrigins) == 0 || (len(cfg.AllowedOrigins) == 1 && cfg.AllowedOrigins[0] == "*") {
		r.Use(cors.Default())
	} else {
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowOrigins = cfg.AllowedOrigins
		corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
		corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
		r.Use(cors.New(corsConfig))
	}

	r.GET("/health", healthHandler.Health)

	api := r.Group("/api/v1")

	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.GET("/google/login", authHandler.GoogleLogin)
	auth.GET("/google/callback", authHandler.GoogleCallback)

	users := api.Group("/users")
	users.Use(middleware.AuthMiddleware(authService))
	users.GET("", userHandler.List)
	users.GET("/:id", userHandler.Get)
	users.PUT("/:id", userHandler.Update)
	users.DELETE("/:id", userHandler.Delete)

	return r
}
