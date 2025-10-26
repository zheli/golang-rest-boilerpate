package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/example/golang-rest-boilerplate/internal/config"
	"github.com/example/golang-rest-boilerplate/internal/models"
	"github.com/example/golang-rest-boilerplate/internal/repository"
	"github.com/example/golang-rest-boilerplate/internal/service"
)

func setupRepository(t *testing.T) *repository.UserRepository {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}))
	return repository.NewUserRepository(db)
}

func TestRegisterAndLogin(t *testing.T) {
	repo := setupRepository(t)
	cfg := &config.Config{JWTSecret: "secret", JWTIssuer: "test", TokenExpireMinutes: 60}
	authService := service.NewAuthService(repo, cfg)

	user, err := authService.Register(context.Background(), "Alice", "alice@example.com", "Password123")
	require.NoError(t, err)
	require.Equal(t, "Alice", user.Name)
	require.NotEmpty(t, user.PasswordHash)

	token, loggedInUser, err := authService.Login(context.Background(), "alice@example.com", "Password123")
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Equal(t, user.Email, loggedInUser.Email)

	claims, err := authService.ParseToken(token)
	require.NoError(t, err)
	require.Equal(t, user.Email, claims.Email)
}
