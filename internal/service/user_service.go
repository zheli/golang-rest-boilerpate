package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/golang-rest-boilerplate/internal/models"
	"github.com/example/golang-rest-boilerplate/internal/repository"
)

// UserService contains business logic for user management.
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService constructs a new UserService.
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// List retrieves all users.
func (s *UserService) List(ctx context.Context) ([]models.User, error) {
	return s.repo.List(ctx)
}

// Get retrieves a user by ID.
func (s *UserService) Get(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

// Update updates the user's name.
func (s *UserService) Update(ctx context.Context, id uuid.UUID, name string) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Name = name
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Delete removes a user by ID.
func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
