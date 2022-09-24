package services

import (
	"alterra-agmc-day-5-6/internal/models"
	"alterra-agmc-day-5-6/internal/repositories"
	"alterra-agmc-day-5-6/pkg/jwt"
	"context"
)

type UserService interface {
	Login(ctx context.Context, email string, password string) (string, error)
	FindAll(ctx context.Context) ([]*models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.User, userID uint) (*models.User, error)
	DeleteByID(ctx context.Context, id uint, userID uint) error
}

type userServiceImpl struct {
	userRepository repositories.UserRepository
}

// Create implements UserService
func (s *userServiceImpl) Create(ctx context.Context, user *models.User) (*models.User, error) {
	return s.userRepository.Create(ctx, user)
}

// DeleteByID implements UserService
func (s *userServiceImpl) DeleteByID(ctx context.Context, id uint, userID uint) error {
	if id != userID {
		return ErrUnauthorized{}
	}
	return s.userRepository.DeleteByID(ctx, id)
}

// FindAll implements UserService
func (s *userServiceImpl) FindAll(ctx context.Context) ([]*models.User, error) {
	return s.userRepository.FindAll(ctx)
}

// FindByID implements UserService
func (s *userServiceImpl) FindByID(ctx context.Context, id uint) (*models.User, error) {
	return s.userRepository.FindByID(ctx, id)
}

// Update implements UserService
func (s *userServiceImpl) Update(ctx context.Context, user *models.User, userID uint) (*models.User, error) {
	if user.ID != userID {
		return nil, ErrUnauthorized{}
	}
	return s.userRepository.Update(ctx, user)
}

// Login implements UserService
func (s *userServiceImpl) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return "", ErrUnauthorized{}
	}
	if user.Password != password {
		return "", ErrUnauthorized{}
	}
	token, err := jwt.NewToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func NewUserService(
	userRepository repositories.UserRepository,
) UserService {
	return &userServiceImpl{
		userRepository: userRepository,
	}
}
