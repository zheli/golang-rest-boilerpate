package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/example/golang-rest-boilerplate/internal/config"
	"github.com/example/golang-rest-boilerplate/internal/models"
	"github.com/example/golang-rest-boilerplate/internal/repository"
)

// ErrInvalidCredentials represents invalid login attempts.
var ErrInvalidCredentials = errors.New("invalid credentials")

// AuthService handles authentication-related operations.
type AuthService struct {
	repo              *repository.UserRepository
	jwtSecret         []byte
	jwtIssuer         string
	tokenExpirePeriod time.Duration
}

// Claims represents JWT claims structure.
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// NewAuthService creates a new AuthService.
func NewAuthService(repo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		repo:              repo,
		jwtSecret:         []byte(cfg.JWTSecret),
		jwtIssuer:         cfg.JWTIssuer,
		tokenExpirePeriod: time.Duration(cfg.TokenExpireMinutes) * time.Minute,
	}
}

// Register creates a new user with hashed password.
func (s *AuthService) Register(ctx context.Context, name, email, password string) (*models.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(passwordHash),
		Provider:     "local",
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user using email and password.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, *models.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// FindOrCreateOAuthUser creates a user account based on OAuth provider information.
func (s *AuthService) FindOrCreateOAuthUser(ctx context.Context, name, email, provider, providerID string) (*models.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err == nil {
		// existing user, update provider info if empty
		if user.Provider == "" {
			user.Provider = provider
			user.ProviderID = providerID
			if err := s.repo.Update(ctx, user); err != nil {
				return nil, err
			}
		}
		return user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user = &models.User{
		Name:       name,
		Email:      email,
		Provider:   provider,
		ProviderID: providerID,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// GenerateToken creates a JWT for the supplied user.
func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.jwtIssuer,
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpirePeriod)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ParseToken validates a JWT and returns its claims.
func (s *AuthService) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrInvalidCredentials
}
